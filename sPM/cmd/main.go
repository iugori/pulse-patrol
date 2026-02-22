package main

import (
	"log"
	"net/http"
	"sPM/internal/greeting"
)

func main() {
	// API Route
	http.HandleFunc("GET /hello", greeting.GreetHandler)

	port := ":8081"

	log.Printf("Greeting Service starting on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
