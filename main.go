package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/thewhitewizard/thegraph-mcp-server/internal/mcp"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/chain"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/thegraph"
)

const (
	cancelTimeOut = 5 * time.Second
)

func main() {
	// Define flags
	useSSE := flag.Bool("sse", false, "Use SSE server mode (default is stdin/stdout)")
	port := flag.String("port", "", "Port for SSE server (defaults to PORT env var or 4000)")
	theGraphURL := flag.String("thegraph-url", "", "TheGraph URL, default "+thegraph.DEFAULT_URL)
	chainRPC := flag.String("rpc", "", "RPC for chain interaction, default "+chain.DEFAULT_URL)

	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	useSSEEnv := getEnv("USE_SSE", "")
	if useSSEEnv == "true" {
		*useSSE = true
	} else if useSSEEnv == "false" {
		*useSSE = false
	}

	// If chainRPC flag not set, get from env or use default
	if *chainRPC == "" {
		*chainRPC = getEnv("CHAIN_RPC", chain.DEFAULT_URL)
	}

	// If theGraphURL flag not set, get from env or use default
	if *theGraphURL == "" {
		*theGraphURL = getEnv("THEGRAPH_URL", thegraph.DEFAULT_URL)
	}

	// If port flag not set, get from env or use default
	if *port == "" {
		*port = getEnv("PORT", "4000")
	}

	logLevel := getEnv("LOG_LEVEL", "info")

	// Configure logging
	switch logLevel {
	case "debug":
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	default:
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// Initialize clients
	thegraphCient := thegraph.NewClient(*theGraphURL)
	chainClient := chain.NewClient(*chainRPC)

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"TheGraph MCP Server",
		"1.0.0",
	)

	// Register tools
	mcp.RegisterTools(mcpServer, thegraphCient, chainClient)

	if *useSSE {
		// SSE server mode
		log.Printf("Starting in SSE mode...")
		runSSEServer(mcpServer, *port)
	} else {
		// Default StdIO server mode
		log.Printf("Starting in StdIO mode...")
		runStdIOServer(mcpServer)
	}
}

func runSSEServer(mcpServer *server.MCPServer, port string) {
	// Create custom SSE server
	sseServer := mcp.NewCustomSSEServer(mcpServer)

	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting TheGraph MCP Server in SSE mode on %s", addr)

		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interruption signal
	<-ctx.Done()
	stop()
	log.Println("Shutting down server...")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cancelTimeOut)
	defer cancel()

	// Attempt graceful shutdown
	if err := sseServer.GracefulShutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func runStdIOServer(mcpServer *server.MCPServer) {
	// Create custom StdIO server
	stdioServer := mcp.NewCustomStdioServer(mcpServer)

	// Configure error logger to write to stderr
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	stdioServer.SetErrorLogger(errorLogger)

	log.Printf("Starting TheGraph MCP Server in StdIO mode")

	// Start listening on stdin/stdout
	if err := stdioServer.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
