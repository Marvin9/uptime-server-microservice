package utils

import "testing"

func TestEncryption(t *testing.T) {
	password := "abc"
	hash, _ := HashAndSalt(password)
	unHash := ComparePassword(hash, "abc")
	if !unHash {
		t.Errorf("For password %v, does not match after hasing", password)
	}
}
