.PHONY: all build clean proto test run help build-api run-api

BINARY_NAME=auraserver
API_BINARY_NAME=apiserver
PROTO_DIR=pkg/api/v1
GEN_DIR=gen/go/provisioning/v1
PROTOC_BIN=$(HOME)/.local/bin/protoc

all: proto build

help:
	@echo "Aura Build System"
	@echo ""
	@echo "Targets:"
	@echo "  proto        - Generate Go code from protobuf definitions"
	@echo "  build        - Build the auraserver binary"
	@echo "  build-api    - Build the apiserver binary"
	@echo "  run          - Run the provisioning server"
	@echo "  run-api      - Run the API server"
	@echo "  test         - Run tests"
	@echo "  clean        - Remove build artifacts"
	@echo "  all          - Generate proto and build (default)"

proto:
	@echo "Generating protobuf code..."
	@mkdir -p $(GEN_DIR)
	@$(PROTOC_BIN) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) --proto_path=$(HOME)/.local/include \
		$(PROTO_DIR)/provisioning.proto
	@mv provisioning.pb.go provisioning_grpc.pb.go $(GEN_DIR)/
	@echo "Protobuf generation complete"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/auraserver
	@echo "Build complete: bin/$(BINARY_NAME)"

build-api:
	@echo "Building $(API_BINARY_NAME)..."
	@go build -o bin/$(API_BINARY_NAME) ./cmd/apiserver
	@echo "Build complete: bin/$(API_BINARY_NAME)"

run:
	@echo "Starting Aura Provisioning Server..."
	@go run ./cmd/auraserver

run-api:
	@echo "Starting Aura API Server..."
	@go run ./cmd/apiserver

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"