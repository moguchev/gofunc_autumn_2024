# gofunc_autumn_2024

## Инициализация репозитория

* `go mod init github.com/moguchev/gofunc_autumn_2024`
____

## Генерация кода из protobuf схем 

### Создание protobuf схемы

* `mkdir -p proto/api/example`
* Создаем `service.proto` и `message.proto` в `./proto/api/example/v1`:
    + `touch ./proto/api/example/v1/service.proto`
    + `touch ./proto/api/example/v1/message.proto`
* Описываем protobuf схему нашего сервиса

### Установка и конфигурация Buf

* Устанавливаем **Buf**
    + `go install github.com/bufbuild/buf/cmd/buf@v1.36.0`
    + `buf --version`
* Инициализируем конфиг **Buf** 
    + `buf config init` - создастся файл `buf.yaml` 
* Указываем путь до наших _protobuf_ фалйов в [buf.yaml](./buf.yaml):
```yaml
modules:
  - path: proto
```
* Создадим конфигурацию `buf.gen.yaml` для генерации кода
    + `touch buf.gen.yaml`
    + Указываем необходимые плагины для генерации go кода (см [buf.gen.yaml](./buf.gen.yaml))

### Установка зависимостей
#### Устанавливаем необходимые плагины gRPC

* `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` - для генерации структур из gRPC сообщений
* `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` - для генерации gRPC сервера и клиента
* `go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest` - для генерации gRPC-Gateway proxy
* `go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest` - для генерации swagger спецификации

#### Устанавливаем protobuf зависимости (googleapis, protovalidate, ...)

##### Вариант 1. Buf Schema Registry (BSR) (вы не в России или у вас VPN)

* Указыаемыем зависимости из **Buf Schema Registry** (BSR) в `buf.yaml`
```yaml
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
  - buf.build/bufbuild/protovalidate
```

* `buf dep update` - создаст файл `buf.lock`

##### Вариант 2. Ручной вендоринг зависимостей (вы в России и у вас нет VPN)

* Вручную (или с помощью скриптов) скачиваем необходимые protobuf файлы и складываем в директорию `vendor.proto`
    - `make vendor` - в нашем примере написан _Makefile_ `vendor.proto.mk`, в котором происходит вендоринг необходимых protobuf файлов
* Указываем путь до скаченных protobuf файлов в `buf.yaml`
```yaml
modules:
  - path: proto
  - path: vendor.protobuf
```

### Генерация go файлов

* `buf build` - проверяем, что все необходимые зависимости корректно установлены
* `buf genereate` - генерация golang кода

### Форматирование protobuf

*  `buf format -w` - запуск форматирования

#### Автоматическое форматирование в VSCode

* Установите расширение [Buf](https://marketplace.visualstudio.com/items?itemName=bufbuild.vscode-buf)
* Выберете Buf в качестве форматера по умолчанию: `⌘ + Shift + P` -> `Fromat Document` -> `Configure Default Formatter`
* Добавьте следующие настройки в `settings.json` файл для форматирования protobuf файлов при сохранении:
```json
"[proto]": {
  "editor.formatOnSave": true
}
```

### Линтер protobuf

* Настраиваем правила для линтера
    - [Linting Overview](https://buf.build/docs/lint/overview)
    - [Rules and Categories](https://buf.build/docs/lint/rules)
    - В `buf.yaml` указываем необходимые нам настройки
    ```yaml
    lint:
      use:
        - STANDARD
        - COMMENTS
    ```
    - В случае вендоринга сторонних protobuf файлов вручную, стоит указать их в разделе `ignore` (не всегда работает)
    ```yaml
    lint:
      use:
        - STANDARD
        - COMMENTS
      ignore:
        - vendor.protobuf
    ```
* `buf lint -v` - запуск линтера


## rk-boot

### rk-grpc

1. `go get github.com/rookie-ninja/rk-grpc/v2` - устанавливаем rk-grpc
2. Создайте в корне проекта файл `boot.yaml` - конфигурация для rk-boot
```yaml
grpc:
  - name: example # Название нашего entry
    description: "example server"
    enabled: true # Можем отключить entry при необходимости (мало ли ...)
    port: 82 # Порт на котором будем принимать входящие gRPC запросы
    enableReflection: true # Включить gRPC reflection (в основном нужно для grpcurl, grpcli, postman)
    gwPort: 80 # Порт на котором будем принимать входящие HTTP запросы в gRPC-Gateway
    gwOption: # Настройки опций gRPC-Gateway
      marshal: # https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson#MarshalOptions
        multiline: false
        emitUnpopulated: true
        indent: ""
        allowPartial: false
        useProtoNames: false
        useEnumNumbers: false
      unmarshal: # https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson#UnmarshalOptions
        allowPartial: false
        discardUnknown: true
```
3. Создайте в корне проекта файл `config.go`
```go
package config

import (
	_ "embed"
)

//go:embed boot.yaml
var Boot []byte
```
4. Создайте `main.go` в директории `cmd/example`:
```go
package main

import (
	"context"
	"log"

	config "github.com/moguchev/gofunc_autumn_2024"
	examplev1 "github.com/moguchev/gofunc_autumn_2024/internal/app/api/example/v1"
	pb "github.com/moguchev/gofunc_autumn_2024/pkg/api/example/v1"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация нашего RPC обработчика
	srv, err := examplev1.NewExampleServiceServerImplementation()
	if err != nil {
		log.Fatalf("couldn't create server: %v", err)
	}

	// Загрузжаем entries из конфигурации (boot.yaml).
	boot := rkboot.NewBoot(
		rkboot.WithBootConfigRaw(config.Boot),
	)

	// Получение GrpcEntry
	grpcEntry := rkgrpc.GetGrpcEntry("example") // название entry
	// Регистрация gRPC сервера
	grpcEntry.AddRegFuncGrpc(func(server *grpc.Server) { pb.RegisterExampleServiceServer(server, srv) })
	// Регистрация gRPC-Gateway proxy
	grpcEntry.AddRegFuncGw(pb.RegisterExampleServiceHandlerFromEndpoint)

	// Bootstrap entry
	boot.Bootstrap(ctx)

	// Ждем сигнала выключения
	boot.WaitForShutdownSig(ctx)
}
```
5. `go run ./cmd/example` - запуск сервиса
```
2024-09-15T15:57:03.400+0300    INFO    boot/grpc_entry.go:1060 Bootstrap grpcEntry     {"eventId": "c60f8189-d910-491b-b938-4ac59cfe422c", "entryName": "example", "entryType": "gRPCEntry"}
2024-09-15T15:57:03.400+0300    INFO    boot/grpc_entry.go:761  gRPC_port:82
2024-09-15T15:57:03.400+0300    INFO    boot/grpc_entry.go:762  gateway_port:80
```
### Документация API

1. Добавим в `boot.yaml` у нашего gRPC entry `example` следующий конфиг:
```yaml
    commonService: # Swagger UI клиент для RK сервиса
      enabled: true
    sw: # Swagger UI клиент: https://github.com/swagger-api/swagger-ui
      enabled: true
      path: "swagger"
      jsonPaths:
        - swagger
      headers: []
    docs: # Встроенный экземпляр RapiDoc https://github.com/rapi-doc/RapiDoc, который можно использовать вместо Swagger
      enabled: true
      path: "docs"
      specPath: "swagger"
      headers: []
      style:
        theme: "light"
      debug: true
```
2. `go run ./cmd/example` - запуск сервиса
```
2024-09-15T16:08:42.262+0300    INFO    boot/grpc_entry.go:1060 Bootstrap grpcEntry     {"eventId": "a636bbf6-afd0-4b06-80e8-0f0f87951e44", "entryName": "example", "entryType": "gRPCEntry"}
2024-09-15T16:08:42.262+0300    INFO    boot/grpc_entry.go:761  gRPC_port:82
2024-09-15T16:08:42.262+0300    INFO    boot/grpc_entry.go:762  gateway_port:80
2024-09-15T16:08:42.262+0300    INFO    boot/grpc_entry.go:765  SwaggerEntry: http://localhost:80/swagger/
2024-09-15T16:08:42.262+0300    INFO    boot/grpc_entry.go:768  DocsEntry: http://localhost:80/docs/
2024-09-15T16:08:42.262+0300    INFO    boot/grpc_entry.go:783  CommonSreviceEntry: http://localhost:80/rk/v1/ready, http://localhost:80/rk/v1/alive, http://localhost:80/rk/v1/info
```

### Observability

#### Метрики

1. Добавим в `boot.yaml` у нашего gRPC entry `example` следующий конфиг:
```yaml
    prom: # Prometheus client will automatically register into grpc-gateway instance at /metrics.
      enabled: true
    middleware:
      prom: # метрики
        enabled: true
```
2. `go run ./cmd/example` - запуск сервиса
3. http://localhost/metrics - endpoint для сбора метрик

#### Профилирование

1. Добавим в `boot.yaml` у нашего gRPC entry `example` следующий конфиг:
```yaml
    pprof: # Профилирование
      enabled: true
      path: "/pprof"
```
2. `go run ./cmd/example` - запуск сервиса
3. http://localhost/pprof - endpoint pprof

#### Логирование

1. Добавим в `boot.yaml` у нашего gRPC entry `example` следующий конфиг для логирования запросов:
```yaml
    middleware:
      logging: # логирование https://github.com/rookie-ninja/rk-query
        enabled: true
        ignore: [""]
        loggerEncoding: "json" # console, json, flatten
        loggerOutputPaths: ["stdout"]
        eventEncoding: "json" # console, json, flatten
        eventOutputPaths: ["stdout"]
```
2. Добавим в `boot.yaml` у нашего gRPC entry `example` следующий конфиг для настройки глобального логирования:
```yaml
logger: # logging with [rk-logger](https://github.com/rookie-ninja/rk-logger)
  - name: zap
    description: "ZAP"
    domain: "*"
    default: true
    zap:
      level: debug
      development: true
      disableCaller: false
      disableStacktrace: true
      encoding: json # console, json, flatten
      outputPaths: ["stdout"]
      errorOutputPaths: ["stderr"]
      encoderConfig:
        timeKey: "ts"
        levelKey: "level"
        nameKey: "logger"
        callerKey: "caller"
        messageKey: "msg"
        stacktraceKey: "stacktrace"
        skipLineEnding: false
        lineEnding: "\n"
        consoleSeparator: "\t"
      sampling: # Optional, default: nil
        initial: 0
        thereafter: 0
      initialFields:
        key: value
event: # logging of RPC with [rk-query](https://github.com/rookie-ninja/rk-query)
  - name: event
    description: "event entry"
    domain: "*"
    encoding: json # console, json, flatten
    default: true
    outputPaths: ["stdout"]
```

#### Трейсинг

1. Добавим в `boot.yaml` у нашего gRPC entry `example` следующий конфиг:
```yaml
    middleware:
      trace: # Трейсинг
        enabled: true
        ignore: [""]
        exporter:
          file:
            enabled: true
            outputPath: "logs/trace.log"
          jaeger:
            agent:
              enabled: false
              host: "localhost"
              port: 6831
            collector:
              enabled: false
              endpoint: "http://localhost:14268/api/traces"
              username: ""
              password: ""
```
