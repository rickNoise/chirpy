package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
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

func (cfg *ApiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
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
	if len(params.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp cannot have an empty body", nil)
		return
	}

	// check if chirp body requires censoring (still valid)
	_, censoredBody := censorChirp(params.Body)

	// If the chirp is valid, you should save it in the database
	parsedUserId, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "provided user_id is not a valid UUID", err)
		return
	}
	dbChirp, err := cfg.DbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   censoredBody,
		UserID: parsedUserId,
	})
	if err != nil {
		if checkForForeignKeyConstraintViolationPostgresql(err) {
			respondWithError(w, http.StatusBadRequest, "could not add chirp to db, likely request contained a user_id that does not exist", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not add chirp to database", err)
		}
		return
	}

	// If creating the record succeeds, respond with a 201 status code and the full chirp resource
	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
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

// Add a GET /api/chirps endpoint that returns all chirps in the database. It should return them in the same structure as the POST /api/chirps endpoint, but as an array. Use a 200 status code for success. Order them by created_at in ascending order.
func (cfg *ApiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DbQueries.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get chirps", err)
		return
	}

	var jsonChirps []Chirp
	for _, dbChirp := range dbChirps {
		jsonChirps = append(jsonChirps, Chirp{
			Id:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, jsonChirps)
}

/* HELPER FUNCTIONS */

// censorChirp accepts a string and checks it for any words needing to be censored.
// It returns a bool that indicates if any censoring was required, and a string of the censored body (empty string if no censor required).

// returns true if passed error is a postgres unique constraint violation error
func checkForUniqueConstraintViolationPostgresql(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

// returns true if passed error is a postgres FK constraint violation
func checkForForeignKeyConstraintViolationPostgresql(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23503"
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

// msg is returned to the requester; err is logged internally
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
