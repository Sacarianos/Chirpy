package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Sacarians/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	JWTSecret      string
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

	// Get DB URL
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Create a new Queries object
	dbQueries := database.New(db)

	//  load the secret from the .env file
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		JWTSecret:      jwtSecret,
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a fileserver thet serves files from the current directory
	fs := http.FileServer(http.Dir(filepathRoot))
	fsWithMetrics := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs))
	// Use the fileserver to handle requests to the root path
	mux.Handle("/app/", fsWithMetrics)

	// API endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetAllChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirp)
	mux.HandleFunc("GET /api/chirps/", apiCfg.handleGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
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
