package main

import (
	"fmt"
	"log"
	"net/http"
)

/* CONSTANTS */
const port = "8080"
const filepathRoot = "."

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("created new http.Server\n")

	fmt.Printf("listening on Addr: %v...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
