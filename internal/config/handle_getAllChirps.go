package config

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"sort"

	"github.com/google/uuid"
	"github.com/rickNoise/chirpy/internal/database"
)

// Update the GET /api/chirps endpoint. It should accept an optional query parameter called author_id.
// If the author_id query parameter is provided, the endpoint should return only the chirps for that author.
// If the author_id query parameter is not provided, the endpoint should return all chirps as it did before.
// For example:
//
// GET http://localhost:8080/api/chirps?author_id=1
//
// Continue sorting the chirps by created_at in ascending order.
func (cfg *ApiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// check for optional author_id query parameter.
	rawAuthorId := r.URL.Query().Get("author_id")

	// check for optional sort query parameter.
	// if the provided value is not in an expected format, fall back on the default.
	acceptableSortValues := []string{"asc", "desc"}
	sortDirection := "asc"
	requestedSort := r.URL.Query().Get("sort")
	if slices.Contains(acceptableSortValues, requestedSort) {
		sortDirection = requestedSort
	}

	// If the author_id query parameter is not provided, the endpoint should return all chirps as it did before.
	// If the author_id query parameter is provided, the endpoint should return only the chirps for that author.
	var dbChirps []database.Chirp
	var dbErr error
	if rawAuthorId == "" {
		dbChirps, dbErr = cfg.DbQueries.GetAllChirps(r.Context())
		if dbErr != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps", dbErr)
			return
		}
	} else {
		parsedAuthorID, err := uuid.Parse(rawAuthorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author_id provided", err)
			return
		}
		dbChirps, dbErr = cfg.DbQueries.GetAllChirpsByAuthorUserId(r.Context(), parsedAuthorID)
		if dbErr != nil {
			if errors.Is(dbErr, sql.ErrNoRows) {
				// case for no chirps found under the provided author_id
				dbChirps = []database.Chirp{}
			} else {
				respondWithError(w, http.StatusInternalServerError, "could not get chirps", err)
				return
			}
		}
	}

	// sort dbChirps based on sort variable; if "asc" do nothing, as db query does this by default
	if sortDirection == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt) })
	}

	// assemble json response by looping over the db chirps
	// chirps should already be sorted by the db queries ASC by created_at
	var jsonChirps []Chirp
	for _, dbChirp := range dbChirps {
		jsonChirps = append(jsonChirps, DatabaseChirpToAPIChirp(dbChirp))
	}

	respondWithJSON(w, http.StatusOK, jsonChirps)
}
