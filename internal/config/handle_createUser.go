package config

import (
	"encoding/json"
	"net/http"

	"github.com/rickNoise/chirpy/internal/auth"
	"github.com/rickNoise/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	if !validatePassword(params.Password) {
		respondWithError(w, http.StatusBadRequest, "invalid password provided", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create user", err)
		return
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		Hashedpassword: hashedPassword,
	})
	if err != nil {
		if checkForUniqueConstraintViolationPostgresql(err) {
			respondWithError(w, http.StatusConflict, "user already exists", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not create user", err)
		}
		return
	}

	respondWithJSON(w, 201, DatabaseUserToAPIUser(dbUser))
}
