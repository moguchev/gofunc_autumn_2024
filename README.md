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