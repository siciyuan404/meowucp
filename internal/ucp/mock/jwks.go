package mock

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
)

type JWK struct {
	KTY string `json:"kty"`
	CRV string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	KID string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func BuildJWK(key *ecdsa.PrivateKey, kid string) JWK {
	return JWK{
		KTY: "EC",
		CRV: "P-256",
		X:   base64.RawURLEncoding.EncodeToString(key.PublicKey.X.Bytes()),
		Y:   base64.RawURLEncoding.EncodeToString(key.PublicKey.Y.Bytes()),
		KID: kid,
		Use: "sig",
		Alg: "ES256",
	}
}
