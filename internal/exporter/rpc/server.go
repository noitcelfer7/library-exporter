package rpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	library_proto "github.com/noitcelfer7/library-proto/gen/go/proto/library"

	"library_exporter/internal/exporter/config"
)

type Server struct {
	library_proto.UnimplementedDataExchangeServiceServer
}

func (server *Server) Exchange(ctx context.Context, req *library_proto.ExchangeRequest) (*library_proto.ExchangeResponse, error) {
	log.Printf("Exchange: %v\n", req)

	return &library_proto.ExchangeResponse{IsSuccessful: true}, nil
}

func Serve(config *config.Config) {
	addr := net.JoinHostPort(config.Grpc.Server.Host, config.Grpc.Server.Port)

	listener, _ := net.Listen("tcp", addr)

	s := grpc.NewServer()

	library_proto.RegisterDataExchangeServiceServer(s, &Server{})

	s.Serve(listener)
}
