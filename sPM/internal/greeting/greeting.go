package greeting

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetWelcomeMessage processes the core business logic for greetings
func GetWelcomeMessage(ctx context.Context, name string) string {
	_, span := otel.Tracer("sPM/greeting").Start(ctx, "GetWelcomeMessage")
	defer span.End()

	span.SetAttributes(attribute.String("name", name))

	if name == "" {
		span.SetAttributes(attribute.Bool("default_greeting", true))
		return "Hello World!"
	}

	span.SetAttributes(attribute.Bool("default_greeting", false))
	return fmt.Sprintf("Hello %s!", name)
}
