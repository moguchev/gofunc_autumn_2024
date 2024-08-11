package main

import (
	"context"
	"log"

	"buf.build/go/protovalidate"
	config "github.com/moguchev/balun_meetup_2025"
	examplev1 "github.com/moguchev/balun_meetup_2025/internal/app/api/example/v1"
	"github.com/moguchev/balun_meetup_2025/internal/middleware"
	"github.com/moguchev/balun_meetup_2025/pkg/core"
	"github.com/moguchev/balun_meetup_2025/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// protovalidate validator
	validator, err := protovalidate.New(
		protovalidate.WithFailFast(),
	)
	if err != nil {
		log.Fatalf("failed to initialize validator: %v", err)
	}

	// Приложение
	boot := core.NewBoot(config.Boot,
		core.WithUnaryInterceptors(
			middleware.WithProtovalidateUnaryServerInterceptor(validator),
		),
	)

	// Конфиг приложения
	cfg := boot.Config()

	logger.Debug("read config", // Для примера
		zap.String("key", cfg.GetString("key")),
		zap.String("key1", cfg.GetString("key1")),
		zap.String("key2", cfg.GetString("key2")),
	)

	// Инициализация нашего RPC обработчика
	srv, err := examplev1.NewExampleServiceServerImplementation(validator)
	if err != nil {
		log.Fatalf("couldn't create server: %v", err)
	}

	// Ждем сигнала выключения
	boot.Run(ctx, srv)
}
