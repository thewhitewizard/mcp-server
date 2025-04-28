.PHONY: build clean install

# Default target
all: build

# Build the server
build:
	go build -o bin/thegraph-mcp-server .

# Clean build artifacts
clean:
	rm -f bin/thegraph-mcp-server