package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rickNoise/chirpy/internal/auth"
	"github.com/rickNoise/chirpy/internal/database"
)

/* CONSTANTS */
// Access tokens should expire in 1 hour.
const ACCESS_TOKEN_EXPIRATION = time.Hour

// Refresh tokens should expire after 60 days. Expiration time is stored in the database.
const REFRESH_TOKEN_EXPIRATION = time.Hour * 24 * 60

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {

	// decode request
	// expecting {
	//   "password": "04234",
	//   "email": "lane@example.com"
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

	// create access token
	accesTokenExpiration := ACCESS_TOKEN_EXPIRATION
	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.JWTSecret,
		accesTokenExpiration,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not login user", fmt.Errorf("failed to create JWT token: %w", err))
		return
	}

	// create refresh token and store it in the db
	var refreshToken string
	refreshToken, _ = auth.MakeRefreshToken()
	_, err = cfg.DbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(REFRESH_TOKEN_EXPIRATION), // calculate expiration timestamp
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not login user", fmt.Errorf("error storing refresh token in db: %w", err))
		return
	}

	type LoginResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	jsonLoginResponse := LoginResponse{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		IsChirpyRed:  dbUser.IsChirpyRed.Bool,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	// otherwise 200 OK and a copy of the user resource without password
	respondWithJSON(w, http.StatusOK, jsonLoginResponse)
}
