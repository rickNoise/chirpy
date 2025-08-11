package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		},
	}

	// Create a new token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Use token.SignedString to sign the token with the secret key.
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	fmt.Printf("signed string token: %s, err: %s", ss, err)
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Use the jwt.ParseWithClaims function to validate the signature of the JWT and extract the claims into a *jwt.Token struct. An error will be returned if the token is invalid or has expired.

	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not parse token: %w", err)
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		// retrieve user ID stored as a string in the Subject field
		userIdString := claims.Subject
		// return user ID as a uuid.UUID
		id, err := uuid.Parse(userIdString)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("invalid user id: %w", err)
		}
		return id, nil

	} else {
		return uuid.UUID{}, errors.New("unknown claims type, cannot proceed")
	}
}
