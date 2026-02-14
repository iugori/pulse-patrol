# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pulse Patrol is a healthcare SaaS system for real-time patient monitoring and data portability. The system bridges medical IoT devices with clinical management, providing real-time alerts and patient record access.

## Development Commands

### Running Services

Each service is a separate Go module located in its own directory:

```bash
# Run Patient Management Service (REST API on :8081)
cd sPM && go run cmd/main.go

# Run Telemetry & Alerting Service (gRPC on :50051)
cd sTA && go run cmd/main.go
```

### Testing

```bash
# Run tests for sPM service
cd sPM && ./go-test.sh

# Run tests for sTA service
cd sTA && ./go-test.sh

# Or manually with go test
cd sPM && go clean -testcache && go test ./internal/test/... -v
cd sTA && go clean -testcache && go test ./internal/test/... -v
```

### Protobuf Generation (sTA service only)

The sTA service uses gRPC and requires protobuf code generation:

```bash
cd sTA && ./go-protoc.sh
```

This generates Go code from `api/greeting.proto` into `api/greeting.pb.go` and `api/greeting_grpc.pb.go`.

### Dependency Management

```bash
# For sTA service (installs protoc plugins and debugging tools)
cd sTA && ./go-tidy-up.sh

# Manual dependency tidying
cd sPM && go mod tidy
cd sTA && go mod tidy
```

## Architecture

### Service Structure

Both services follow the standard Go project layout:

```
sPM/  or  sTA/
├── cmd/
│   └── main.go          # Service entry point
├── api/
│   ├── handlers.go      # HTTP handlers (sPM only)
│   └── *.proto          # Protocol definitions (sTA only)
├── internal/
│   ├── greeting/        # Business logic
│   └── test/            # Tests
└── pkg/                 # Shared packages (future use)
```

### Microservices Overview

The system is composed of several specialized services documented in `docs/hw6/ARD-v6.1-rfc.md`:

**Core Services (Implemented)**:
- **sPM** (Patient Management): Handles medical records, admissions, and patient transfers via REST API
- **sTA** (Telemetry & Alerting): Processes real-time medical equipment data and triggers alerts via gRPC

**Planned Services**:
- **sAAA** (Compliance & Identity): Authentication, authorization, and audit logging
- **sGW** (Integration Gateway): Protocol translation for legacy systems and medical equipment
- **uiWP** (Web Portal): Patient and administrator interface
- **uiCD** (Clinical Dashboard): Clinical staff monitoring interface

### Communication Patterns

As defined in `docs/hw6/adr-v3.1-rfc.md`, the system uses hybrid communication:

- **Synchronous (REST/gRPC)**: User-facing operations requiring immediate feedback
- **Asynchronous (Message Queues)**: High-volume telemetry and audit logging
- **Reactive (WebSockets)**: Real-time alerts pushed to clinical staff

### Key Architectural Constraints

From `docs/hw6/adr-v3.1-rfc.md`:
- **Database Schema Isolation**: Services use exclusive database schemas; cross-service data exchange happens via defined communication edges only
- **Zero-Trust Identity**: All requests must include valid identity tokens (OIDC/OAuth2)
- **MQTT QoS Levels**: Medical equipment uses QoS 2 (exactly once) for alerts, QoS 0 (at most once) for high-frequency waveform data
- **Non-Blocking Audit**: Audit logging is fire-and-forget to prevent blocking clinical operations
- **Data Encryption**: All PII must be encrypted at rest and in transit (HIPAA/GDPR compliance)

## Domain Model

The system implements Domain-Driven Design with bounded contexts documented in `docs/hw6/ARD-v6.1-rfc.md`:

- **Care Coordination & Admissions**: Patient admission lifecycle and inter-company transfers
- **Clinical Records**: Medical history, lab results, and patient records
- **Vital Signs & Monitoring**: Real-time telemetry processing and threshold evaluation
- **Notification & Alerting**: Alert lifecycle from detection to staff acknowledgment
- **Security & Audit**: Access control and compliance logging

## AWS Deployment (Planned)

Target infrastructure per `docs/hw6/ARD-v6.1-rfc.md`:
- Compute: AWS Fargate (ECS) for 24/7 service availability
- Data: Amazon Aurora (PostgreSQL) for records, Amazon Timestream for telemetry
- Messaging: SNS/SQS for async communication, AWS IoT Core for MQTT
- Real-time: AWS AppSync for WebSocket connections
