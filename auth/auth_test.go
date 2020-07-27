package auth

import "testing"

func TestGetToken(t *testing.T) {
	secret := "4567890kjhgfdwertyucvbmnbvdfgu098765edcvbnyds234rvbhuuhko0opo1qasxc"
	a := New(secret)
	tc := TokenClaims{
		ID:    "1",
		Email: "abuarquemf@gmail.com",
	}
	_, err := a.GetToken(&tc)
	if err != nil {
		t.Errorf("failed to create token, get error %q", err)
	}
}

func TestIsValid(t *testing.T) {
	secret := "4567890kjhgfdwertyucvbmnbvdfgu098765edcvbnyds234rvbhuuhko0opo1qasxc"
	a := New(secret)
	tc := TokenClaims{
		ID:    "1",
		Email: "abuarquemf@gmail.com",
	}
	encrypted, err := a.GetToken(&tc)
	if err != nil {
		t.Errorf("failed to create token, get error %q", err)
	}
	res, err := a.IsValid(encrypted)
	if err != nil {
		t.Errorf("failed to check if is valid with error %q", err)
	}
	if res == false {
		t.Errorf("want %v, got %v", true, res)
	}
}
