package main

import (
	"context"
	"log"
	"net"

	library_proto "github.com/noitcelfer7/library-proto/gen/go/proto/library"
	"google.golang.org/grpc"
)

type server struct {
	library_proto.UnimplementedDataExchangeServiceServer
}

func (s *server) Exchange(ctx context.Context, req *library_proto.ExchangeRequest) (*library_proto.ExchangeResponse, error) {
	log.Printf("Received: %+v", req)
	return &library_proto.ExchangeResponse{IsSuccessful: true}, nil
}

func main() {
	lis, _ := net.Listen("tcp", "0.0.0.0:12345")
	s := grpc.NewServer()
	library_proto.RegisterDataExchangeServiceServer(s, &server{})
	s.Serve(lis)
}
