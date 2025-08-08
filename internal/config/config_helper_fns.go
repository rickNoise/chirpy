package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

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
