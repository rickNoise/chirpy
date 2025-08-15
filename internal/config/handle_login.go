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

/* CONSTANTS */
const ACCESS_TOKEN_EXPIRATION = time.Hour

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {

	// decode request
	// expecting {
	//   "password": "04234",
	//   "email": "lane@example.com",
	// }
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	accesTokenExpiration := ACCESS_TOKEN_EXPIRATION

	// create token
	createdToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.JWTSecret,
		accesTokenExpiration,
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
