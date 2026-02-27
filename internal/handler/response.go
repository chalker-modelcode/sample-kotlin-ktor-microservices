package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON encodes v as JSON and writes it to the response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// errorResponse is the standard JSON error body.
type errorResponse struct {
	Error string `json:"error"`
}

// writeError writes a JSON error response with the given status code and message.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
