include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"
GOBIN=$(LOCAL_BIN)

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/protovalidate-go

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	mkdir -p pkg/swagger
	make generate-user-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.png,*.json'

generate-user-api:
	mkdir -p pkg/user_v1
	protoc \
	--proto_path vendor.protogen \
	--proto_path api/user_v1 \
	--proto_path vendor.protogen/google \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/user_v1/*.proto

build:
	GOOS=linux GOARCH=amd64 go build -o auth cmd/main.go

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

start:
	docker-compose up -d pg-local migrator-local
	go run cmd/main.go

restart:
	go run cmd/main.go

vendor-proto:
	@if [ ! -d vendor.protogen/buf/validate/priv ]; then \
		mkdir -p vendor.protogen/buf/validate && \
		curl -sSL https://raw.githubusercontent.com/bufbuild/protovalidate/refs/heads/main/proto/protovalidate/buf/validate/validate.proto -o vendor.protogen/buf/validate/validate.proto ;\
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone --depth 1 https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi