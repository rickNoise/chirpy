package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rickNoise/chirpy/internal/auth"
)

const MAX_AUTH_EXPIRATION_DURATION_IN_SECONDS = 3600

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {

	// decode request
	// expecting {
	//   "password": "04234",
	//   "email": "lane@example.com",
	//   "expires_in_seconds": "3600" (optional)
	// }
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode request body", err)
		return
	}

	// check for user record
	dbUser, err := cfg.DbQueries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// check password matches the record
	// return 401 Unauthorised if no match
	if err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var tokenExpiration time.Duration
	// handle case of no user-provided value
	if params.ExpiresInSeconds == nil {
		tokenExpiration = getValidExpirationTimeInSeconds(0)
	} else {
		tokenExpiration = getValidExpirationTimeInSeconds(*params.ExpiresInSeconds)
	}

	// create token
	createdToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.JWTSecret,
		tokenExpiration,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not login user", fmt.Errorf("failed to create JWT token: %w", err))
		return
	}

	type LoginResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	jsonLoginResponse := LoginResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     createdToken,
	}
	// otherwise 200 OK and a copy of the user resource without password
	respondWithJSON(w, http.StatusOK, jsonLoginResponse)
}

// max auth expiration time is a constant
func getValidExpirationTimeInSeconds(requestedExpiration int) time.Duration {
	if requestedExpiration <= 0 || requestedExpiration > MAX_AUTH_EXPIRATION_DURATION_IN_SECONDS {
		return MAX_AUTH_EXPIRATION_DURATION_IN_SECONDS * time.Second
	} else {
		return time.Duration(requestedExpiration) * time.Second
	}
}
