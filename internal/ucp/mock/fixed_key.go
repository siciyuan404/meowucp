package mock

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

func FixedKey(kid string) (*ecdsa.PrivateKey, JWK) {
	curve := elliptic.P256()
	priv := &ecdsa.PrivateKey{D: big.NewInt(1)}
	priv.PublicKey.Curve = curve
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(priv.D.Bytes())
	return priv, BuildJWK(priv, kid)
}
