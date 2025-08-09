package config

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rickNoise/chirpy/internal/auth"
)

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

	jsonUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	// otherwise 200 OK and a copy of the user resource without password
	respondWithJSON(w, http.StatusOK, jsonUser)

}
