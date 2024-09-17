# Путь до завендореных protobuf файлов
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

# vendor
vendor:	.vendor-reset .vendor-googleapis .vendor-google-protobuf .vendor-protovalidate .vendor-protoc-gen-openapiv2 .vendor-tidy

# delete VENDOR_PROTO_PATH
.vendor-reset:
	rm -rf $(VENDOR_PROTO_PATH)
	mkdir -p $(VENDOR_PROTO_PATH)

# Устанавливаем proto описания google/protobuf
.vendor-google-protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_PROTO_PATH)/protobuf &&\
	cd $(VENDOR_PROTO_PATH)/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/google/protobuf
	find $(VENDOR_PROTO_PATH)/protobuf/src/google/protobuf -maxdepth 1 -type f -exec mv {} $(VENDOR_PROTO_PATH)/google/protobuf \;
	rm -rf $(VENDOR_PROTO_PATH)/protobuf

# Устанавливаем proto описания google/api
.vendor-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/googleapis &&\
	cd $(VENDOR_PROTO_PATH)/googleapis &&\
	git sparse-checkout set --no-cone google/api &&\
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/google/api
	find $(VENDOR_PROTO_PATH)/googleapis/google/api -maxdepth 1 -type f -exec mv {} $(VENDOR_PROTO_PATH)/google/api \;
	rm -rf $(VENDOR_PROTO_PATH)/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
.vendor-protoc-gen-openapiv2:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_PATH)/grpc-gateway && \
 	cd $(VENDOR_PROTO_PATH)/grpc-gateway && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_PATH)/grpc-gateway/protoc-gen-openapiv2/options $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_PATH)/grpc-gateway

# Устанавливаем proto описания buf/validate/validate.proto
.vendor-protovalidate:
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate $(VENDOR_PROTO_PATH)/protovalidate && \
	cd $(VENDOR_PROTO_PATH)/protovalidate
	git checkout
	mv $(VENDOR_PROTO_PATH)/protovalidate/proto/protovalidate/buf $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/protovalidate

# delete all non .proto files
.vendor-tidy:
	find $(VENDOR_PROTO_PATH) -type f ! -name "*.proto" -delete
	find $(VENDOR_PROTO_PATH) -type f \( -name "*unittest*.proto" -o -name "*test*.proto" -o -name "*sample*.proto" \) -delete
	find $(VENDOR_PROTO_PATH) -empty -type d -delete

# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.vendor-reset \
	.vendor-google-protobuf \
	.vendor-googleapis \
	.vendor-protoc-gen-openapiv2 \
	.vendor-protovalidate \
	.vendor-tidy \
	vendor