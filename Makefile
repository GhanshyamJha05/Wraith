.PHONY: all build test eval demo clean

BINARY_NAME=wraith
BIN_DIR=bin

all: build

build:
	@echo "==> Building $(BINARY_NAME)..."
	go build -trimpath -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/wraith

test:
	@echo "==> Running tests..."
	go test -v ./...

eval:
	@echo "==> Running evaluation harness..."
	go run ./eval/main.go

demo: build
	@echo "==> Running demo..."
	./$(BIN_DIR)/$(BINARY_NAME) demo --debug

clean:
	@echo "==> Cleaning..."
	rm -rf $(BIN_DIR)
