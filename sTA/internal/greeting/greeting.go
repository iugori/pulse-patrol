package greeting

import (
	"context"
	"sTA/api"
)

type Server struct {
	api.UnimplementedGreeterServer
}

func (s *Server) SayHello(_ context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	return &api.HelloResponse{
		Message: "Hello " + in.GetName(),
	}, nil
}
