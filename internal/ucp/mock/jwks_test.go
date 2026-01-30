package mock

import "testing"

func TestBuildJWK(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	jwk := BuildJWK(key, "kid_1")
	if jwk.KID != "kid_1" {
		t.Fatalf("expected kid to be set")
	}
	if jwk.X == "" || jwk.Y == "" {
		t.Fatalf("expected x/y to be set")
	}
}
