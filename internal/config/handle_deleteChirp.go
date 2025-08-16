package config

import (
	"net/http"

	"github.com/google/uuid"
)

// Add a new DELETE /api/chirps/{chirpID} route to your server that deletes a chirp from the database by its id.
// This is an authenticated endpoint, so be sure to check the token in the header. Only allow the deletion of a chirp if the user is the author of the chirp.
// If they are not, return a 403 status code.
// If the chirp is deleted successfully, return a 204 status code.
// If the chirp is not found, return a 404 status code.
func (cfg *ApiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// authenticate requesting user
	requestingUserID, ok := cfg.authenticateUser(w, r)
	if !ok {
		return // helper already wrote the error response
	}

	// parse request param as a uuid
	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}

	// check requesting user id matches the author of the chirp to delete
	chirpToDelete, err := cfg.DbQueries.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}
	if requestingUserID != chirpToDelete.UserID {
		respondWithError(w, http.StatusForbidden, "", nil)
		return
	}

	_, err = cfg.DbQueries.DeleteChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete chirp", err)
		return
	}

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}
