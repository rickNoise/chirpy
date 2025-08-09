package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash the password using the bcrypt.GenerateFromPassword function. Bcrypt is a secure hash function that is intended for use with passwords.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}

	// confirm hash worked properly
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return "", fmt.Errorf("problem when confirming password hashing worked: %w", err)
	}

	// safe to cast []byte to string here
	return string(hashedPassword), err
}

// Use the bcrypt.CompareHashAndPassword function to compare the password that the user entered in the HTTP request with the password that is stored in the database.
func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("password did not match: %w", err)
	}
	return nil
}
