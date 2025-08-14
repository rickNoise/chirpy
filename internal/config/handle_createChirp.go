package config

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rickNoise/chirpy/internal/auth"
	"github.com/rickNoise/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding req json body", err)
		return
	}

	// check request for valid Authorization header
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid bearer token in request", err)
		return
	}

	// determine posting user by JWT
	parsedUserId, err := auth.ValidateJWT(bearerToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user bearer token", err)
		return
	}

	// check length of chirp body
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	if len(params.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp cannot have an empty body", nil)
		return
	}

	// check if chirp body requires censoring (still valid)
	_, censoredBody := censorChirp(params.Body)

	dbChirp, err := cfg.DbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   censoredBody,
		UserID: parsedUserId,
	})
	if err != nil {
		if checkForForeignKeyConstraintViolationPostgresql(err) {
			respondWithError(w, http.StatusBadRequest, "could not add chirp to db, likely request contained a user_id that does not exist", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not add chirp to database", err)
		}
		return
	}

	// If creating the record succeeds, respond with a 201 status code and the full chirp resource
	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	})
}
