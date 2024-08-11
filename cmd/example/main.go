package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	examplev1 "github.com/moguchev/gofunc_autumn_2024/internal/app/api/example/v1"
	pb "github.com/moguchev/gofunc_autumn_2024/pkg/api/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const GRPC_PORT = ":80"

func main() {
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// TLS/SSL
	tlsConfig := &tls.Config{ /* ... */ }

	srv, err := examplev1.NewExampleServiceServerImplementation()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC-Gateway
	mux := runtime.NewServeMux()
	err = pb.RegisterExampleServiceHandlerServer(context.TODO(), mux, srv)
	if err != nil {
		log.Fatalf("failed to register handler:: %v", err)
	}
	httpServer := &http.Server{
		Handler:   mux,
		Addr:      GRPC_PORT,
		TLSConfig: tlsConfig,
	}
	go httpServer.ListenAndServe()

	// gRPC
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	)
	pb.RegisterExampleServiceServer(grpcServer, srv)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
