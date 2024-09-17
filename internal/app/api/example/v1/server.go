package examplev1

import (
	"github.com/bufbuild/protovalidate-go"
	pb "github.com/moguchev/gofunc_autumn_2024/pkg/api/example/v1"
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
