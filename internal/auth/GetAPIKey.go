package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Extract the api key from the Authorization header, which is expected to be in this format:
// Authorization: ApiKey THE_KEY_HERE
//
// Expects to find an "Authorization" header in headers.
// Expects a value in the format: "ApiKey <KEY>".
// Returns <KEY>. Does not guarantee that <KEY> is a valid api key.
func GetAPIKey(headers http.Header) (string, error) {
	// Check if an Authorization header exists
	rawAuthHeaderValue := headers.Get("Authorization")
	if rawAuthHeaderValue == "" {
		return "", errors.New("cannot get value, input does not have an Authorization header")
	}

	// expecting the header value of the form "ApiKey <KEY>"
	parsedAuthHeaderValue := strings.Fields(rawAuthHeaderValue)

	if len(parsedAuthHeaderValue) != 2 {
		return "", fmt.Errorf("cannot get value, unexpected Header format: %v", parsedAuthHeaderValue)
	}

	if parsedAuthHeaderValue[0] != "ApiKey" {
		return "", fmt.Errorf("cannot get value, malformed API key: %v", parsedAuthHeaderValue)
	}

	// we want to return <KEY> only
	return parsedAuthHeaderValue[1], nil
}
