package examplev1

import (
	"context"

	pb "github.com/moguchev/gofunc_autumn_2024/pkg/api/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateNote - метод создания заметки
func (i *ExampleServiceServerImplementation) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	// Валидация (Вариант 1)
	if err := i.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Ваша реализация
	// time.Sleep(100 * time.Millisecond)

	return &pb.CreateNoteResponse{NoteId: "какой-то id"}, nil
}
