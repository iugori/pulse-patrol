package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sPM/internal/greeting"
	"sPM/internal/telemetry"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Initialize OpenTelemetry
	signozEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	shutdown, err := telemetry.InitTracer("sPM", signozEndpoint)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	// Create ServeMux with OpenTelemetry instrumentation
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", greeting.GreetHandler)

	// Wrap with OpenTelemetry middleware
	handler := otelhttp.NewHandler(mux, "sPM")

	port := ":8081"
	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Greeting Service starting on %s...", port)
		log.Printf("Sending traces to SigNoz at %s", signozEndpoint)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
