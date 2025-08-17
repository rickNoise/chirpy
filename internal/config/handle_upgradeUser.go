package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// Add a POST /api/polka/webhooks endpoint. It should accept a request of this shape:
//
//	{
//	  "event": "user.upgraded",
//	  "data": {
//	    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
//	  }
//	}
//
// If the event is anything other than user.upgraded, the endpoint should immediately respond with a 204 status code - we don't care about any other events.
// If the event is user.upgraded, then it should update the user in the database, and mark that they are a Chirpy Red member.
// If the user is upgraded successfully, the endpoint should respond with a 204 status code and an empty response body. If the user can't be found, the endpoint should respond with a 404 status code.
func (cfg *ApiConfig) HandleUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode request body", err)
		return
	}

	// If the event is anything other than user.upgraded, the endpoint should immediately respond with a 204 status code - we don't care about any other events.
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DbQueries.UpgradeUserToChirpyRedById(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows found - this is a 404 case
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		// Any other database error - this is a 500 case
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	// if user is upgraded successfully, respond with 204 No Content and an empty response body
	w.WriteHeader(http.StatusNoContent)
}
