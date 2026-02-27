package handler

import "net/http"

// HealthCheck handles GET /health requests, returning a simple liveness response.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"UP"}`))
}
