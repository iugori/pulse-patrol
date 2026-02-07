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
