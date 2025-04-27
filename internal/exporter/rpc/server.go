package rpc

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	library_proto "github.com/noitcelfer7/library-proto/gen/go/proto/library"

	"library_exporter/internal/exporter/config"
	postgresql "library_exporter/internal/exporter/database"
	"library_exporter/internal/exporter/schema"
)

type Server struct {
	db *postgresql.Database

	library_proto.UnimplementedDataExchangeServiceServer
}

func (server *Server) Exchange(ctx context.Context, req *library_proto.ExchangeRequest) (*library_proto.ExchangeResponse, error) {
	log.Printf("Exchange: %v\n", req)

	server.db.InsertRecord(schema.Record{
		AuthorFirstName: req.AuthorFirstName,
		AuthorLastName:  req.AuthorLastName,

		BookIsbn:  req.BookIsbn,
		BookTitle: req.BookTitle,

		GenreTitle: req.GenreTitle,

		IssueDate:       req.IssueDate,
		IssuePeriod:     req.IssuePeriod,
		IssueReturnDate: sql.NullString{String: req.GetIssueReturnDate()},

		ReaderFirstName:   req.ReaderFirstName,
		ReaderLastName:    req.ReaderLastName,
		ReaderPhoneNumber: req.ReaderPhoneNumber,
	})

	return &library_proto.ExchangeResponse{IsSuccessful: true}, nil
}

func Serve(config *config.Config, db *postgresql.Database) {
	cert, err := tls.LoadX509KeyPair("server-cert.pem", "server-key.pem")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}

	creds := credentials.NewTLS(tlsConfig)

	if err != nil {
		log.Fatal(err)
	}

	addr := net.JoinHostPort(config.Grpc.Server.Host, config.Grpc.Server.Port)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		panic(fmt.Sprintf("net.Listen Error: %v", err))
	}

	s := grpc.NewServer(grpc.Creds((creds)))

	library_proto.RegisterDataExchangeServiceServer(s, &Server{db: db})

	s.Serve(listener)
}
