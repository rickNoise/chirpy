package auth

import (
	"testing"
)

// TestHashPassword calls HashPassword to check if a string is properly hashed.
func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "testPassword123!"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("error when using HashPassword: %v", err)
	}

	if CheckPasswordHash(password, hashedPassword) != nil {
		t.Errorf("CheckPasswordHash claiming the password and hashed password don't match, but they should!: %v", err)
	}
}
