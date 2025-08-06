package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

const maxChirpLength = 140

type ApiConfig struct {
	fileserverHits atomic.Int32
}

type chirp struct {
	Body string `json:"body"`
}

type jsonError struct {
	Error string `json:"error"`
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) ReadinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	template := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>	
`
	hits := fmt.Sprintf(template, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hits))
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirpBody := chirp{}
	err := decoder.Decode(&chirpBody)
	if err != nil {
		marshalledErr, innerErr := json.Marshal(jsonError{Error: err.Error()})
		if innerErr != nil {
			fmt.Printf("error marshalling error into json: %s\n", innerErr)
			w.WriteHeader(500)
			return
		}
		fmt.Printf("error decoding chirp: %s", err)
		w.WriteHeader(400)
		w.Write(marshalledErr)
		return
	}

	// check length of chirp body
	if len(chirpBody.Body) > maxChirpLength {
		marshalledErr, err := json.Marshal(jsonError{Error: "Chirp is too long"})
		if err != nil {
			fmt.Printf("error marshalling error into json: %s\n", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(marshalledErr)
		return
	}

	// chirp is valid
	successJSON := struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	}
	dat, err := json.Marshal(successJSON)
	if err != nil {
		fmt.Printf("error marshalling valid response json: %v\n", successJSON)
	}
	w.WriteHeader(200)
	w.Write(dat)
}
