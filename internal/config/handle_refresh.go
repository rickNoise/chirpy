package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// // Create a POST /api/refresh endpoint. This new endpoint does not accept a request body, but does require a refresh token to be present in the headers, in the same Authorization: Bearer <token> format.
// Look up the token in the database. If it doesn't exist, or if it's expired, respond with a 401 status code. Otherwise, respond with a 200 code and this shape:
//
//	{
//	  "token": "<token>"
//	}
//
// The token field should be a newly created access token for the given user that expires in 1 hour. I wrote a GetUserFromRefreshToken SQL query.
func (cfg *ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// first ensure valid refresh token is provided in request heaader
	rawToken := r.Header.Get("Authorization")
	parsedToken := strings.Fields(rawToken)
	// catches if header value not in format: "Bearer <token>"
	if len(parsedToken) != 2 {
		respondWithError(w, http.StatusUnauthorized, "cannot refresh token", fmt.Errorf("cannot get Bearer Token, unexpected Header format: %s", rawToken))
		return
	}
	// assumes header value format: "Bearer <token>"
	tokenString := parsedToken[1]
	dbRefreshToken, err := cfg.DbQueries.GetRefreshTokenByTokenString(context.Background(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	type RefreshResponse struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, RefreshResponse{
		Token: dbRefreshToken.Token,
	})
}
