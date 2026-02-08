# Architecture Decision Record - v5.1

> **<span style="font-size: 1.3em">Pulse Patrol</span>**
>
> *Develop a software system for healthcare that collects and manages patient data,
> integrates with medical equipment, provides web access for patients and authorized personnel,
> alerts staff for abnormal values, and supports patient transfers between healthcare providers.*

<!-- TOC -->
* [Architecture Decision Record - v5.1](#architecture-decision-record---v51)
  * [TL;DR](#tldr)
  * [1. Context](#1-context)
    * [Scope](#scope)
    * [Architectural Nodes](#architectural-nodes)
  * [3. Decision](#3-decision)
    * [Project structure](#project-structure)
    * [REST service implementation](#rest-service-implementation)
      * [sPM/api/handlers.go](#spmapihandlersgo)
      * [sPM/cmd/main.go](#spmcmdmaingo)
      * [sPM/internal/greeting/greeting.go](#spminternalgreetinggreetinggo)
    * [gRPC service implementation](#grpc-service-implementation)
      * [sTA/api/greeting.proto](#staapigreetingproto)
      * [sTA/cmd/main.go](#stacmdmaingo)
      * [sTA/internal/greeting/greeting.go](#stainternalgreetinggreetinggo)
<!-- TOC -->

## TL;DR

The source code is available here:

https://github.com/iugori/pulse-patrol

## 1. Context

### Scope

[//]: # (<< What is the problem we are trying to solve >>)

> * Based on your chosen case study, select a service to implement in Golang.
    For that service, write a basic REST API and run it locally.
    The code should be pushed to GitHub.
> * Based on your chosen case study, select a service to implement in Golang.
    For that service, write a basic gRPC API and run it locally.
    The code should be pushed to your personal GitHub.

### Architectural Nodes

This brief system overview is reproduced from ADR v3.1 (Homework 3)

**Containers**

- **sPM - Patient Management Services**: Core logic for medical records, admission forms, and inter-company transfers.
- **sTA - Telemetry & Alerting Services**: Processes real-time data from medical equipment and triggers notifications
  for abnormal values.

## 3. Decision

[//]: # (<< What alternative was chosen and why >>)

Implementation highlights

* I provided basic "Hello World!" service implementations in Golang that can be accessed via REST and gRPC.
* Each project structure reproduces the conventional folder structure for Go projects as recommended during the training
  session.
* Additionally shell scripts wer provided to automate special project configuration and test automation.

The project source code is available at https://github.com/iugori/pulse-patrol but in case there are access issues here
are the most important points:

### Project structure

<p align="center">
  <img src="misc/h5-project-structure.png" width="200">
</p>

### REST service implementation

#### sPM/api/handlers.go

```go
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sPM/internal/greeting"
)

// GreetHandler translates HTTP requests into service calls
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	// Business logic call
	message := greeting.GetWelcomeMessage(name)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})

	if err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		return
	}
}
```

#### sPM/cmd/main.go

```go
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

```

#### sPM/internal/greeting/greeting.go

```go
package greeting

import "fmt"

// GetWelcomeMessage processes the core business logic for greetings
func GetWelcomeMessage(name string) string {
	if name == "" {
		return "Hello World!"
	}
	return fmt.Sprintf("Hello %s!", name)
}
```

### gRPC service implementation

#### sTA/api/greeting.proto

```protobuf
syntax = "proto3";

option go_package = "sTA/api";

package greeting;

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

#### sTA/cmd/main.go

```go
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"sTA/api"
	"sTA/internal/greeting"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterGreeterServer(s, &greeting.Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

#### sTA/internal/greeting/greeting.go

```go
package greeting

import (
	"context"
	"fmt"
	"sTA/api"
)

type Server struct {
	api.UnimplementedGreeterServer
}

func (s *Server) SayHello(_ context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	name := in.GetName()
	var responseMessage string

	if name == "" {
		responseMessage = "Hello World!"
	} else {
		responseMessage = fmt.Sprintf("Hello %s!", name)
	}

	return &api.HelloResponse{
		Message: responseMessage,
	}, nil
}
```
