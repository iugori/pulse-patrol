# OpenTelemetry Integration

The sPM service is instrumented with OpenTelemetry to send distributed traces to SigNoz.

## What Was Added

1. **Telemetry Package** (`internal/telemetry/otel.go`):
   - Initializes OpenTelemetry tracer provider
   - Configures OTLP gRPC exporter for SigNoz
   - Sets up W3C Trace Context propagation
   - Handles graceful shutdown

2. **HTTP Instrumentation**:
   - Automatic HTTP request/response tracing via `otelhttp` middleware
   - Captures HTTP method, path, status code, and duration

3. **Custom Spans**:
   - `GreetHandler`: Traces HTTP handler execution with query parameters
   - `GetWelcomeMessage`: Traces business logic with input attributes

4. **Graceful Shutdown**:
   - Proper signal handling (SIGINT, SIGTERM)
   - Flushes pending telemetry data before exit

## Running with SigNoz

### Prerequisites

Make sure SigNoz is running locally (https://github.com/SigNoz). If using Docker:

```bash
brew install --cask docker

git clone -b main https://github.com/SigNoz/signoz.git && cd signoz/deploy/

./install.sh
```

SigNoz will be available at:
- UI: http://localhost:8080
- OTLP gRPC Receiver: localhost:4317

### Running the Service

By default, the service sends traces to `localhost:4317`. To start:

```bash
cd sPM
go run cmd/main.go
```

To use a different SigNoz endpoint:

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=my-signoz-host:4317 go run cmd/main.go
```

### Testing the Instrumentation

Make some requests to generate traces:

```bash
# Basic greeting
curl http://localhost:8081/hello

# Personalized greeting
curl "http://localhost:8081/hello?name=Alice"
```

### Viewing Traces in SigNoz

1. Open SigNoz UI at http://localhost:8080
2. Navigate to "Services" to see the `sPM` service
3. Click on the service to view traces
4. Explore individual traces to see:
   - HTTP request details (method, path, status)
   - Handler execution time
   - Business logic spans with attributes
   - Query parameters and custom attributes

## Trace Structure

Each HTTP request creates a trace with the following span hierarchy:

```
GET /hello (otelhttp instrumentation)
└── GreetHandler (custom span)
    └── GetWelcomeMessage (custom span)
```

### Span Attributes

- **HTTP Span** (auto-instrumented):
  - `http.method`: HTTP method
  - `http.route`: Route pattern
  - `http.status_code`: Response status
  - `http.target`: Request path with query string

- **GreetHandler**:
  - `query.name`: Name query parameter
  - `response.success`: Whether response was encoded successfully

- **GetWelcomeMessage**:
  - `name`: Input name parameter
  - `default_greeting`: Whether default greeting was used

## Configuration

The service accepts the following environment variable:

- `OTEL_EXPORTER_OTLP_ENDPOINT`: SigNoz collector endpoint (default: `localhost:4317`)

## Dependencies

OpenTelemetry dependencies added to `go.mod`:

- `go.opentelemetry.io/otel`: Core OpenTelemetry API
- `go.opentelemetry.io/otel/sdk`: OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`: OTLP gRPC exporter
- `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`: HTTP instrumentation
- `google.golang.org/grpc`: gRPC for OTLP transport

## Performance testing 

```bash
brew install wrk

wrk -t4 -c100 -d20s --latency http://localhost:8081
```