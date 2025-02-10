# Variables
PROJECT_NAME = reminder-svc
GATEWAY_NAME = gateway-svc
GO = go
PROTO_DIR = internal/proto
PROTO_FILE = $(PROTO_DIR)/reminder.proto
PROTO_OUT = $(PROTO_DIR)
PORT = 50052

# Targets
.PHONY: build run proto clean

# Build the service
build:
	@echo "Building $(PROJECT_NAME)..."
	$(GO) build -o bin/$(PROJECT_NAME) ./cmd/main.go

# Run the service
run: build
	@echo "Running $(PROJECT_NAME) on port $(PORT)..."
	./bin/$(PROJECT_NAME)

# Generate Go code from .proto file
proto:
	@echo "Generating Go code from $(PROTO_FILE)..."
	protoc --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_FILE)
	cp $(PROTO_DIR)/reminder.pb.go ../$(GATEWAY_NAME)/internal/proto/reminder.pb.go
	cp $(PROTO_DIR)/reminder_grpc.pb.go ../$(GATEWAY_NAME)/internal/proto/reminder_grpc.pb.go

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/$(PROJECT_NAME)