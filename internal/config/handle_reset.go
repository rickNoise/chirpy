package config

import "net/http"

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "cannot reset all user data in a non-dev environment", nil)
		return
	}

	err := cfg.DbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete user data", err)
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
