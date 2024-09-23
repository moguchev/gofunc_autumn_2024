package examplev1

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/moguchev/gofunc_autumn_2024/pkg/api/example/v1"
	"google.golang.org/grpc"
)

// ExampleServiceServerImplementation - реализует интерфейс pb.ExampleServiceServer
type ExampleServiceServerImplementation struct {
	// UnimplementedExampleServiceServer - для обеспечения прямой совместимости нашей
	// реализации с интерфейсом pb.ExampleServiceServer.
	pb.UnimplementedExampleServiceServer

	// Другие зависимости ...
	validator *protovalidate.Validator
}

// NewExampleServiceServerImplementation - конструктор ExampleServiceServerImplementation
func NewExampleServiceServerImplementation(validator *protovalidate.Validator) (*ExampleServiceServerImplementation, error) {
	return &ExampleServiceServerImplementation{
		validator: validator,
	}, nil
}

// Для совместимости с core

// Name - returns config name of our entry
func (i *ExampleServiceServerImplementation) Name() string {
	return "example"
}

// GrpcRegFunc - returns grpc register rout
func (i *ExampleServiceServerImplementation) GrpcRegFunc() func(server *grpc.Server) {
	return func(server *grpc.Server) { pb.RegisterExampleServiceServer(server, i) }
}

// GwRegFunc - returns config name of our entry
func (i *ExampleServiceServerImplementation) GwRegFunc() func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error {
	return pb.RegisterExampleServiceHandlerFromEndpoint
}
