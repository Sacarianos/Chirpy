package main

import "net/http"

// Readiness endpoint
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// add content-Type header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	// Write the status code
	w.WriteHeader(http.StatusOK)

	// Write the body message
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
