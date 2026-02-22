package greeting

import (
	"context"
	"fmt"
)

type Server struct {
	UnimplementedGreeterServer
}

func (s *Server) SayHello(_ context.Context, in *HelloRequest) (*HelloResponse, error) {
	name := in.GetName()
	var responseMessage string

	if name == "" {
		responseMessage = "Hello World!"
	} else {
		responseMessage = fmt.Sprintf("Hello %s!", name)
	}

	return &HelloResponse{
		Message: responseMessage,
	}, nil
}
