package config

import (
	"context"
	"net/http"
)

func (cfg *ApiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DbQueries.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get chirps", err)
		return
	}

	var jsonChirps []Chirp
	for _, dbChirp := range dbChirps {
		jsonChirps = append(jsonChirps, DatabaseChirpToAPIChirp(dbChirp))
	}

	respondWithJSON(w, http.StatusOK, jsonChirps)
}
