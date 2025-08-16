package config

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rickNoise/chirpy/internal/auth"
)

// continue with update logic using userID
// authenticateUser is a helper function to handle user authentication on incoming http requests.
// If there is an issue with the request's authorization header, this function will modify the http.ResponseWriter input.
// If this function returns ok=false in the bool output, the calling function should simply return immediately to let the Response be sent.
func (cfg *ApiConfig) authenticateUser(w http.ResponseWriter, r *http.Request) (userID uuid.UUID, ok bool) {
	// check request for valid Authorization header
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid bearer token in request", err)
		return uuid.Nil, false
	}

	// check if the provided token is valid
	userID, err = auth.ValidateJWT(bearerToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token could not be validated", err)
		return uuid.Nil, false
	}

	return userID, true
}
