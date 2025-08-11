package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	tokenSecret := "testSecret"
	userID := uuid.New()

	signedTokenString, err := MakeJWT(userID, tokenSecret, time.Second)
	if err != nil {
		t.Errorf("failed to MakeJWT token: %v", err)
	}

	parsedUserId, err := ValidateJWT(signedTokenString, tokenSecret)
	if err != nil {
		t.Errorf("failed to ValidateJWT: %v", err)
	}

	if parsedUserId != userID {
		t.Errorf("userID parsed from signed token string %v not equal to original userID %v", parsedUserId, userID)
	}
}

func TestMakeJWTAndValidateJWTWithExpiredTimeout(t *testing.T) {
	tokenSecret := "testSecret"
	userID := uuid.New()
	expiresIn := time.Second

	signedTokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to MakeJWT token: %v", err)
	}

	// pause execution to let token expire
	time.Sleep(expiresIn)

	_, err = ValidateJWT(signedTokenString, tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT should have delivered an error because token should be expired!")
	}
}
