package password

import "testing"

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("password123")
	if err != nil {
		t.Fatal(err)
	}
	if !Verify(hash, "password123") {
		t.Fatal("expected password to match")
	}
	if Verify(hash, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}
