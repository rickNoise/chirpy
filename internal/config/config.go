package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/lib/pq"
	"github.com/rickNoise/chirpy/internal/database"
)

const maxChirpLength = 140

type ApiConfig struct {
	fileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
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

func (cfg *ApiConfig) HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type censoredChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding req json body", err)
		return
	}

	// check length of chirp body
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// check if chirp body requires censoring (still valid)
	_, censoredBody := censorChirp(params.Body)

	respondWithJSON(w, http.StatusOK, censoredChirp{
		CleanedBody: censoredBody,
	})
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode request body", err)
		return
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		if checkForUniqueConstraintViolationPostgresql(err) {
			respondWithError(w, http.StatusConflict, "user already exists", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not create user", err)
		}
		return
	}

	// map database.User --> User struct to control JSON keys
	jsonUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 201, jsonUser)
}

/* HELPER FUNCTIONS */

// censorChirp accepts a string and checks it for any words needing to be censored.
// It returns a bool that indicates if any censoring was required, and a string of the censored body (empty string if no censor required).

// returns true if passed error is a postgres unique constraint violation error
func checkForUniqueConstraintViolationPostgresql(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func censorChirp(body string) (bool, string) {
	/*
		Replace all "profane" words with 4 asterisks: ****.

		Assuming the length validation passed, replace any of the following words in the Chirp with the static 4-character string ****:
		kerfuffle
		sharbert
		fornax

		Be sure to match against uppercase versions of the words as well, but not punctuation. "Sharbert!" does not need to be replaced, we'll consider it a different word due to the exclamation point.
	*/

	// hardcoded list of banned words & suitable replacement
	bannedWordsSlice := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	replacementForBannedWord := "****"

	// create a map simulating a set
	bannedWordMap := make(map[string]struct{})
	for _, word := range bannedWordsSlice {
		bannedWordMap[word] = struct{}{}
	}

	// split into words on space
	bodyWords := strings.Split(body, " ")

	// initialise
	censoredBody := []string{}
	anyCensorDone := false

	// check each word against the banned word set
	for _, word := range bodyWords {
		if _, found := bannedWordMap[strings.ToLower(word)]; found {
			censoredBody = append(censoredBody, replacementForBannedWord)
			anyCensorDone = true
		} else {
			censoredBody = append(censoredBody, word)
		}
	}

	return anyCensorDone, strings.Join(censoredBody, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		fmt.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
