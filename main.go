package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rickNoise/chirpy/internal/config"
	"github.com/rickNoise/chirpy/internal/database"

	_ "github.com/lib/pq"
)

/* CONSTANTS */
const port = "8080"
const filepathRoot = "."

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialise database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set in .env file")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	fmt.Println("successfully connected to db")

	// Create an instance of apiConfig
	apiCfg := &config.ApiConfig{}
	apiCfg.DbQueries = database.New(db)

	mux := http.NewServeMux()

	/* /APP/ PATH PREFIX - SERVE WEBSITE */
	// Create the fileserver
	fileServer := http.FileServer(http.Dir(filepathRoot))
	// Strip the prefix
	strippedHandler := http.StripPrefix("/app", fileServer)
	// Wrap with middleware
	wrappedHandler := apiCfg.MiddlewareMetricsInc(strippedHandler)
	// Handle /app/ pattern
	mux.Handle("/app/", wrappedHandler)

	/* /API/ PATH PREFIX - SERVE API */
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.HandlerValidateChirp)
	mux.HandleFunc("GET /api/healthz", apiCfg.ReadinessHandler)

	/* /ADMIN/ PATH PREFIX */
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
