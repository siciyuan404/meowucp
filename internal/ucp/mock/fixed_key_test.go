package mock

import "testing"

func TestFixedKey(t *testing.T) {
	key, jwk := FixedKey("fixed-kid")
	if key == nil {
		t.Fatalf("expected key")
	}
	if jwk.KID != "fixed-kid" {
		t.Fatalf("expected kid to be set")
	}
}
