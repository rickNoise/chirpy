package config

import (
	"context"
	"fmt"
	"net/http"
)

// Create a new POST /api/revoke endpoint. This new endpoint does not accept a request body, but does require a refresh token to be present in the headers, in the same Authorization: Bearer <token> format.
// Revoke the token in the database that matches the token that was passed in the header of the request by setting the revoked_at to the current timestamp. Remember that any time you update a record, you should also be updating the updated_at timestamp.
// Respond with a 204 status code. A 204 status means the request was successful but no body is returned.
func (cfg *ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getTokenStringFromAuthorizationHeader(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "cannot refresh token", fmt.Errorf("cannot get Bearer Token, unexpected Header formats"))
		return
	}

	_, err = cfg.DbQueries.RevokeRefreshToken(context.Background(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
