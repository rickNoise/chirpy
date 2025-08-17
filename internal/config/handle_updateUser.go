package config

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rickNoise/chirpy/internal/auth"
	"github.com/rickNoise/chirpy/internal/database"
)

// Add a PUT /api/users endpoint so that users can update their own (but not others') email and password. It requires:
// *An access token in the header
// *A new password and email in the request body
// The request must have BOTH an email and a passowrd, they are both required.
func (cfg *ApiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	userID, ok := cfg.authenticateUser(w, r)
	if !ok {
		return // helper already wrote the error response
	}

	if !validatePassword(params.Password) {
		respondWithError(w, http.StatusBadRequest, "invalid password provided", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem updating password", err)
		return
	}

	updatedUser, err := cfg.DbQueries.UpdateEmailAndPasswordByUserId(
		context.Background(),
		database.UpdateEmailAndPasswordByUserIdParams{
			Newemail:          params.Email,
			Newhashedpassword: hashedPassword,
			Userid:            userID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	})
}
