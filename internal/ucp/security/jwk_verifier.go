package security

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	signatureHeader = "UCP-Signature"
	keyIDHeader     = "UCP-Key-Id"
)

type JWKVerifier struct {
	jwkSetURL  string
	clockSkew  time.Duration
	cacheTTL   time.Duration
	mu         sync.RWMutex
	keys       map[string]*ecdsa.PublicKey
	fetchedAt  time.Time
	skipVerify bool
	nonceStore NonceStore
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	KTY string `json:"kty"`
	CRV string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	KID string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
}

func NewJWKVerifier(jwkSetURL string, clockSkewSeconds int) *JWKVerifier {
	return &JWKVerifier{
		jwkSetURL: jwkSetURL,
		clockSkew: time.Duration(clockSkewSeconds) * time.Second,
		cacheTTL:  10 * time.Minute,
		keys:      map[string]*ecdsa.PublicKey{},
	}
}

func (v *JWKVerifier) SetSkipVerify(skip bool) {
	v.skipVerify = skip
}

type NonceStore interface {
	Seen(nonce string) (bool, error)
	Mark(nonce string) error
}

func (v *JWKVerifier) SetNonceStore(store NonceStore) {
	v.nonceStore = store
}

func (v *JWKVerifier) Verify(r *http.Request, body []byte) error {
	if v.skipVerify {
		return nil
	}
	if strings.TrimSpace(v.jwkSetURL) == "" {
		return errors.New("jwk_set_url_not_configured")
	}

	stamp, signature, err := parseSignatureHeader(r.Header.Get(signatureHeader))
	if err != nil {
		return err
	}

	if v.clockSkew > 0 {
		if time.Since(stamp) > v.clockSkew || stamp.Sub(time.Now()) > v.clockSkew {
			return errors.New("signature_timestamp_expired")
		}
	}

	kid := strings.TrimSpace(r.Header.Get(keyIDHeader))
	if kid == "" {
		return errors.New("missing_key_id")
	}

	pubKey, err := v.getKey(r.Context(), kid)
	if err != nil {
		return err
	}

	msg := []byte(fmt.Sprintf("%d.%s", stamp.Unix(), string(body)))
	hash := sha256.Sum256(msg)
	if !ecdsa.VerifyASN1(pubKey, hash[:], signature) {
		return errors.New("invalid_signature")
	}

	if v.nonceStore != nil {
		nonce := buildNonce(r.Header.Get(signatureHeader), kid)
		seen, err := v.nonceStore.Seen(nonce)
		if err != nil {
			return err
		}
		if seen {
			return errors.New("replay_detected")
		}
		if err := v.nonceStore.Mark(nonce); err != nil {
			return err
		}
	}
	return nil
}

func (v *JWKVerifier) getKey(ctx context.Context, kid string) (*ecdsa.PublicKey, error) {
	v.mu.RLock()
	if key, ok := v.keys[kid]; ok && !v.isExpired() {
		v.mu.RUnlock()
		return key, nil
	}
	v.mu.RUnlock()

	v.mu.Lock()
	defer v.mu.Unlock()
	if key, ok := v.keys[kid]; ok && !v.isExpired() {
		return key, nil
	}

	keys, err := fetchJWKKeys(ctx, v.jwkSetURL)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, errors.New("no_keys_available")
	}
	v.keys = keys
	v.fetchedAt = time.Now()
	key, ok := v.keys[kid]
	if !ok {
		return nil, errors.New("key_not_found")
	}
	return key, nil
}

func (v *JWKVerifier) isExpired() bool {
	if v.fetchedAt.IsZero() {
		return true
	}
	return time.Since(v.fetchedAt) > v.cacheTTL
}

func fetchJWKKeys(ctx context.Context, url string) (map[string]*ecdsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("jwk_fetch_failed: %d", resp.StatusCode)
	}

	var set JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, err
	}

	keys := map[string]*ecdsa.PublicKey{}
	for _, jwk := range set.Keys {
		pubKey, err := parseECDSAKey(jwk)
		if err != nil {
			continue
		}
		if jwk.KID != "" {
			keys[jwk.KID] = pubKey
		}
	}

	return keys, nil
}

func parseECDSAKey(jwk JWK) (*ecdsa.PublicKey, error) {
	if jwk.KTY != "EC" || jwk.CRV != "P-256" {
		return nil, errors.New("unsupported_key_type")
	}

	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, err
	}
	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, err
	}

	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}

	return key, nil
}

func parseSignatureHeader(header string) (time.Time, []byte, error) {
	parts := strings.Split(header, ",")
	var ts string
	var sig string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "t=") {
			ts = strings.TrimPrefix(part, "t=")
		}
		if strings.HasPrefix(part, "v1=") {
			sig = strings.TrimPrefix(part, "v1=")
		}
	}
	if ts == "" || sig == "" {
		return time.Time{}, nil, errors.New("missing_signature")
	}

	unix, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, nil, errors.New("invalid_timestamp")
	}

	signature, err := base64.RawURLEncoding.DecodeString(sig)
	if err != nil {
		return time.Time{}, nil, errors.New("invalid_signature_encoding")
	}

	return time.Unix(unix, 0), signature, nil
}

func buildNonce(signatureHeaderValue string, keyID string) string {
	seed := signatureHeaderValue + ":" + keyID
	digest := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(digest[:])
}
