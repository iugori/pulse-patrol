package greeting

import (
	"encoding/json"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("sPM/greeting")

// GreetHandler translates HTTP requests into service calls
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "GreetHandler")
	defer span.End()

	name := r.URL.Query().Get("name")
	span.SetAttributes(attribute.String("query.name", name))

	// Business logic call
	message := GetWelcomeMessage(ctx, name)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})

	if err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		span.RecordError(err)
		return
	}

	span.SetAttributes(attribute.Bool("response.success", true))
}
