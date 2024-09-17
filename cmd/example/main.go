package main

import (
	"context"
	"log"

	"github.com/bufbuild/protovalidate-go"
	config "github.com/moguchev/gofunc_autumn_2024"
	examplev1 "github.com/moguchev/gofunc_autumn_2024/internal/app/api/example/v1"
	"github.com/moguchev/gofunc_autumn_2024/internal/middleware"
	pb "github.com/moguchev/gofunc_autumn_2024/pkg/api/example/v1"
	"github.com/moguchev/gofunc_autumn_2024/pkg/logger"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("start")

	// Загрузжаем basic entries из конфигурации (boot.yaml).
	rkentry.BootstrapBuiltInEntryFromYAML(config.Boot)
	rkentry.BootstrapPluginEntryFromYAML(config.Boot)

	// Загрузжаем entries из конфигурации (boot.yaml).
	boot := rkboot.NewBoot(
		rkboot.WithBootConfigRaw(config.Boot),
	)

	// Конфиг приложения
	cfg := rkentry.GlobalAppCtx.GetConfigEntry("config")

	logger.Debug("read config", // Для примера
		zap.String("key", cfg.GetString("key")),
		zap.String("key1", cfg.GetString("key1")),
		zap.String("key2", cfg.GetString("key2")),
	)

	// protovalidate validator
	validator, err := protovalidate.New(
		protovalidate.WithFailFast(true),
	)
	if err != nil {
		log.Fatalf("failed to initialize validator: %v", err)
	}

	// Инициализация нашего RPC обработчика
	srv, err := examplev1.NewExampleServiceServerImplementation(validator)
	if err != nil {
		log.Fatalf("couldn't create server: %v", err)
	}

	// Получение GrpcEntry
	grpcEntry := rkgrpc.GetGrpcEntry("example") // название entry
	// Регистрация gRPC сервера
	grpcEntry.AddRegFuncGrpc(func(server *grpc.Server) { pb.RegisterExampleServiceServer(server, srv) })
	// Регистрация gRPC-Gateway proxy
	grpcEntry.AddRegFuncGw(pb.RegisterExampleServiceHandlerFromEndpoint)
	// Добавляем наши middleware
	grpcEntry.AddUnaryInterceptors(
		middleware.WithProtovalidateUnaryServerInterceptor(validator),
	)
	grpcEntry.AddStreamInterceptors()

	// Bootstrap entry
	// grpcEntry.Bootstrap(ctx)
	boot.Bootstrap(ctx)

	// Ждем сигнала выключения
	boot.WaitForShutdownSig(ctx)
}
