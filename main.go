package main

import (
	"log"
	"net/http"

	"github.com/rickNoise/chirpy/internal/config"
)

/* CONSTANTS */
const port = "8080"
const filepathRoot = "."

func main() {
	// Create an instance of apiConfig
	apiCfg := &config.ApiConfig{}

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
