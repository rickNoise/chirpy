package config

import (
	"net/http"
	"sync/atomic"

	"github.com/rickNoise/chirpy/internal/database"
)

const maxChirpLength = 140

type ApiConfig struct {
	fileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	JWTSecret      string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
