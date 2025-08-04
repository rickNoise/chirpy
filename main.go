package main

import (
	"fmt"
	"log"
	"net/http"
)

/* CONSTANTS */
const SERVER_ADDR = ":8080"

func main() {
	serveMux := http.NewServeMux()

	s := http.Server{
		Addr:    SERVER_ADDR,
		Handler: serveMux,
	}
	fmt.Printf("created new http.Server\n")

	fmt.Printf("listening on Addr: %v", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("error on ListenAndServe: %v", err)
	}
}
