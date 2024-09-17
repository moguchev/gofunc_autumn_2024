include vendor.proto.mk

# Используем bin в текущей директории для установки зависимостей
LOCAL_BIN := $(CURDIR)/bin

# Устанавливаем необходимые зависимости
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install github.com/bufbuild/buf/cmd/buf@v1.41.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.22.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0

# Импорт внешних proto зависимостей через buf
.buf-deps:
	$(info run buf dep update...)

	PATH="$(LOCAL_BIN):$(PATH)" $(LOCAL_BIN)/buf dep update

# Генерация .pb файлов с помощью buf
.buf-generate:
	$(info run buf generate...)

	PATH="$(LOCAL_BIN):$(PATH)" $(LOCAL_BIN)/buf generate

# Форматирование protobuf файлов
.buf-format:
	$(info run buf format...)

	$(LOCAL_BIN)/buf format -w	

# Генерация .pb файлов
generate: .bin-deps .buf-generate .buf-format

# Генерация .pb файлов без установки плагинов
fast-generate: .buf-generate .buf-format

# Линтер protobuf файлов
.buf-lint:
	$(info run buf lint...)

	$(LOCAL_BIN)/buf lint	

# Линтер
lint: .buf-lint

# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.buf-deps \
	.buf-generate \
	.buf-format \
	.buf-lint \
	generate \
	lint