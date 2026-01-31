.PHONY: build test clean run fmt vet docker-build docker-run

# Build the binary
build:
	go build -o sparkyfitness-mcp ./cmd/sparkyfitness-mcp

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -f sparkyfitness-mcp
	rm -f cmd/sparkyfitness-mcp/sparkyfitness-mcp

# Run the server (requires environment variables)
run: build
	./sparkyfitness-mcp

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Build Docker image
docker-build:
	docker build -t sparkyfitness-mcp .

# Run Docker container (requires environment variables)
docker-run:
	docker run -e SPARKYFITNESS_API_URL -e SPARKYFITNESS_API_KEY sparkyfitness-mcp
