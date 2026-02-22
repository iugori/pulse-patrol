package greeting

import (
	"encoding/json"
	"log"
	"net/http"
)

// GreetHandler translates HTTP requests into service calls
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	// Business logic call
	message := GetWelcomeMessage(name)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})

	if err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		return
	}
}
