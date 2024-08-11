include vendor.proto.mk

# Используем bin в текущей директории для установки зависимостей
LOCAL_BIN := $(CURDIR)/bin

# Путь до buf
BUF_BIN := $(LOCAL_BIN)/buf

# Устанавливаем необходимые зависимости
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install github.com/bufbuild/buf/cmd/buf@v1.56.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.7
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.1

# Устанавливаем необходимые зависимости
.tool-deps: export GOBIN := $(LOCAL_BIN)
.tool-deps:
	$(info Installing dependencies...)

	go get -tool github.com/bufbuild/buf/cmd/buf@v1.56.0
	go get -tool google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.7
	go get -tool google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go get -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.1
	go get -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.1

# Импорт внешних proto зависимостей через buf
.buf-deps:
	$(info run buf dep update...)

	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BIN) dep update

# Генерация .pb файлов с помощью buf
.buf-generate:
	$(info run buf generate...)

	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BIN) generate

# Форматирование protobuf файлов
.buf-format:
	$(info run buf format...)

	$(BUF_BIN) format -w	

# Генерация .pb файлов
generate: .bin-deps .tool-deps .buf-generate .buf-format

# Генерация .pb файлов без установки плагинов
fast-generate: .buf-generate .buf-format

# Линтер protobuf файлов
.buf-lint:
	$(info run buf lint...)

	$(BUF_BIN) lint	

# Линтер
lint: .buf-lint

build:
	$(info building app...)
	go build -o ${LOCAL_BIN} $(CURDIR)/cmd/example


# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.buf-deps \
	.buf-generate \
	.buf-format \
	.buf-lint \
	generate \
	lint \
	build