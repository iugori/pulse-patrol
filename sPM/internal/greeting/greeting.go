package greeting

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetWelcomeMessage processes the core business logic for greetings
func GetWelcomeMessage(ctx context.Context, name string) string {
	_, span := otel.Tracer("sPM/greeting").Start(ctx, "GetWelcomeMessage")
	defer span.End()

	span.SetAttributes(attribute.String("name", name))

	// Simulate processing delay
	delay := simulateProcessingDelay()
	span.SetAttributes(
		attribute.Int("processing.delay_ms", int(delay.Milliseconds())),
		attribute.Bool("processing.slow_path", delay > 500*time.Millisecond),
	)
	time.Sleep(delay)

	if name == "" {
		span.SetAttributes(attribute.Bool("default_greeting", true))
		return "Hello World!"
	}

	span.SetAttributes(attribute.Bool("default_greeting", false))
	return fmt.Sprintf("Hello %s!", name)
}

// simulateProcessingDelay adds realistic variable latency
// - Normal case: 0-500ms random delay
// - 3% probability: 1000-2000ms slow path
func simulateProcessingDelay() time.Duration {
	// 3% chance of slow path
	if rand.IntN(100) < 3 {
		// Slow path: 1000-2000ms
		return time.Duration(1000+rand.IntN(1001)) * time.Millisecond
	}

	// Normal path: 0-500ms
	return time.Duration(rand.IntN(501)) * time.Millisecond
}
