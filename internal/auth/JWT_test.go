package auth

import (
	"fmt"
	"net/http"
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

func TestGetBearerTokenInvalidInput(t *testing.T) {
	// test an invalid call
	invalidHeaders := make(http.Header)
	invalidHeaders.Add("Authorization", "Bearer")

	token, err := GetBearerToken(invalidHeaders)
	if token != "" || err == nil {
		t.Errorf("expected GetBearerToken to fail with invalid input %v, got a token %+v and an err %s", invalidHeaders, token, err)
	}

	// test a valid call
	validHeaders := make(http.Header)
	mockJWTBearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	validHeaders.Add("Authorization", fmt.Sprintf("Bearer %s", mockJWTBearerToken))

	token, err = GetBearerToken(validHeaders)
	if err != nil || token == "" {
		t.Errorf("expected GetBearerToken to succeed with valid input %v", validHeaders)
	}
}
