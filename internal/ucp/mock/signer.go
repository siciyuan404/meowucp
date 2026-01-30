package mock

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func SignPayload(privateKey *ecdsa.PrivateKey, timestamp int64, payload []byte) (string, error) {
	message := buildSignatureMessage(timestamp, payload)
	hash := sha256.Sum256(message)
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(signature), nil
}

func buildSignatureMessage(timestamp int64, payload []byte) []byte {
	return []byte(fmt.Sprintf("%d.%s", timestamp, string(payload)))
}
