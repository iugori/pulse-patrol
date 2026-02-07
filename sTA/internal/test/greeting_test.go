package test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"sTA/api"
	"sTA/internal/greeting"
)

func TestSayHello(t *testing.T) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	// 1. Initialize Server
	s := grpc.NewServer()
	api.RegisterGreeterServer(s, &greeting.Server{})

	// Ensure server stops when test finishes to prevent leaks
	defer s.Stop()

	go func() {
		if err := s.Serve(lis); err != nil {
			// Serve will return an error when s.Stop() is called; that's normal
			return
		}
	}()

	// 2. Initialize Client with modern NewClient API
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 3. Close connection gracefully and handle potential error (Linter fix)
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("failed to close connection: %v", err)
		}
	}()

	// 4. Execute the RPC
	client := api.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &api.HelloRequest{Name: "sTA"})

	// 5. Verify Results
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	if resp.GetMessage() != "Hello sTA!" {
		t.Errorf("Expected 'Hello World', got %s", resp.GetMessage())
	}
}
