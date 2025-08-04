package main

import (
	"log"
	"net/http"

	"github.com/rickNoise/chirpy/handlers"
)

/* CONSTANTS */
const port = "8080"
const filepathRoot = "."

func main() {
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)
	mux.HandleFunc("/healthz", handlers.ReadinessHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
