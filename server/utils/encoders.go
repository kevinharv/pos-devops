package utils

import (
	"fmt"
	"encoding/json"
	"net/http"
)

// JSON encoder for HTTP responses
func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode payload
	err := json.NewEncoder(w).Encode(v); 
	if err != nil {
		return fmt.Errorf("failed to encode JSON with error: %s", err)
	}
	return nil
}

// JSON decoder for HTTP requests
func Decode[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return v, fmt.Errorf("failed to decode JSON with error: %s", err)
	}
	return v, nil
}