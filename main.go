package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1) // Increment the hit counter
		next.ServeHTTP(w, r)      // Call the next handler
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a fileserver thet serves files from the current directory
	fs := http.FileServer(http.Dir(filepathRoot))
	fsWithMetrics := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs))
	// Use the fileserver to handle requests to the root path
	mux.Handle("/app/", fsWithMetrics)

	// API endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// Chirp validation endpoint
	mux.HandleFunc("/api/validate_chirp", apiCfg.handlerValidate)
	// Admin endpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Create a new server struct
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//Start the server
	fmt.Println("Server is running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err) // Log if the server fails to start
	}
}
