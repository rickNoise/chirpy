package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling Get Chirp request for chirpID: %s\n", r.PathValue("chirpID"))

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirpID", nil)
		return
	}

	dbChirp, err := cfg.DbQueries.GetChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "no chirp found with that ID", err)
		return
	}

	respondWithJSON(w, http.StatusOK, DatabaseChirpToAPIChirp(dbChirp))
}
