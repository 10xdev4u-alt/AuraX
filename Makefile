.PHONY: all build clean proto test run help build-api run-api build-ota run-ota build-all docker-build docker-up docker-down

BINARY_NAME=auraserver
API_BINARY_NAME=apiserver
OTA_BINARY_NAME=otaorchestrator
PROTO_DIR=pkg/api/v1
GEN_DIR=gen/go/provisioning/v1
PROTOC_BIN=$(HOME)/.local/bin/protoc

all: proto build

help:
	@echo "AuraX Build System"
	@echo ""
	@echo "Build Targets:"
	@echo "  proto          - Generate Go code from protobuf definitions"
	@echo "  build          - Build the provisioning server binary"
	@echo "  build-api      - Build the API server binary"
	@echo "  build-ota      - Build the OTA orchestrator binary"
	@echo "  build-all      - Build all binaries"
	@echo ""
	@echo "Run Targets:"
	@echo "  run            - Run the provisioning server"
	@echo "  run-api        - Run the API server"
	@echo "  run-ota        - Run the OTA orchestrator"
	@echo ""
	@echo "Docker Targets:"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-up      - Start all services with Docker Compose"
	@echo "  docker-down    - Stop all services"
	@echo ""
	@echo "Utility Targets:"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove build artifacts"
	@echo "  fmt            - Format Go code"
	@echo "  vet            - Run go vet"
	@echo "  all            - Generate proto and build (default)"

proto:
	@echo "Generating protobuf code..."
	@mkdir -p $(GEN_DIR)
	@$(PROTOC_BIN) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) --proto_path=$(HOME)/.local/include \
		$(PROTO_DIR)/provisioning.proto
	@mv provisioning.pb.go provisioning_grpc.pb.go $(GEN_DIR)/
	@echo "✅ Protobuf generation complete"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/auraserver
	@echo "✅ Build complete: bin/$(BINARY_NAME)"

build-api:
	@echo "Building $(API_BINARY_NAME)..."
	@go build -o bin/$(API_BINARY_NAME) ./cmd/apiserver
	@echo "✅ Build complete: bin/$(API_BINARY_NAME)"

build-ota:
	@echo "Building $(OTA_BINARY_NAME)..."
	@go build -o bin/$(OTA_BINARY_NAME) ./cmd/otaorchestrator
	@echo "✅ Build complete: bin/$(OTA_BINARY_NAME)"

build-all: build build-api build-ota
	@echo "✅ All binaries built successfully"

run:
	@echo "Starting Aura Provisioning Server..."
	@go run ./cmd/auraserver

run-api:
	@echo "Starting Aura API Server..."
	@go run ./cmd/apiserver

run-ota:
	@echo "Starting Aura OTA Orchestrator..."
	@go run ./cmd/otaorchestrator

docker-build:
	@echo "Building Docker images..."
	@docker-compose build
	@echo "✅ Docker images built"

docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "✅ Services started"
	@docker-compose ps

docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "✅ Services stopped"

test:
	@echo "Running tests..."
	@go test -v ./...

fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@echo "✅ Code formatted"

vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✅ Go vet passed"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "✅ Clean complete"
