package encryption

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt encrypts a given string
func Encrypt(stringToEncrypt string) (string, error) {
	stringToEncryptAsBytes := []byte(stringToEncrypt)
	hash, err := bcrypt.GenerateFromPassword(stringToEncryptAsBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %q", err)
	}
	return string(hash), nil
}

// Check checks if these two strings are similar
func Check(stringToCheck, stringFromDataSource string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(stringFromDataSource), []byte(stringToCheck))
	if err != nil {
		return false, fmt.Errorf("failed to verify password hash: %q", err)
	}
	return true, nil
}
