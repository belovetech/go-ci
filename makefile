WEB_APP = web
BUILD_DIR = $(PWD)/bin

.PHONY: run_web test bench

web:
	@echo "Building web app..."
	@go build -o $(BUILD_DIR)/$(WEB_APP) cmd/$(WEB_APP)/main.go

run_web: web
	@echo "Running web app..."
	@$(BUILD_DIR)/$(WEB_APP)

test:
	@echo "Running tests..."
	@go test -v ./...

bench:
	@echo "Running benchmarks..."
	@go test -bench=. ./...
