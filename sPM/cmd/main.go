package main

import (
	"log"
	"net/http"
	"sPM/api"
)

func main() {
	// API Route
	http.HandleFunc("GET /hello", api.GreetHandler)

	port := ":8081"

	log.Printf("Greeting Service starting on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
