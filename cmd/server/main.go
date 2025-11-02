package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/aionmcp/aionmcp/internal/core"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ConfigOverrides holds command-line configuration overrides
type ConfigOverrides struct {
	ConfigFile string
	HTTPPort   int
	GRPCPort   int
	LogLevel   string
}

func main() {
	// Define command line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
		configFile  = flag.String("config", "", "Path to configuration file")
		httpPort    = flag.Int("http-port", 0, "HTTP server port (overrides config)")
		grpcPort    = flag.Int("grpc-port", 0, "gRPC server port (overrides config)")
		logLevel    = flag.String("log-level", "", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Println("AionMCP Server v0.1.0")
		fmt.Println("Iteration: 0")
		fmt.Println("Build: development")
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		fmt.Println("AionMCP Server - Autonomous Model Context Protocol Server")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  aionmcp [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Environment Variables:")
		fmt.Println("  AIONMCP_HTTP_PORT     HTTP server port")
		fmt.Println("  AIONMCP_GRPC_PORT     gRPC server port")
		fmt.Println("  AIONMCP_LOG_LEVEL     Log level (debug, info, warn, error)")
		fmt.Println("  AIONMCP_CONFIG        Path to configuration file")
		fmt.Println()
		os.Exit(0)
	}

	// Initialize configuration
	overrides := ConfigOverrides{
		ConfigFile: *configFile,
		HTTPPort:   *httpPort,
		GRPCPort:   *grpcPort,
		LogLevel:   *logLevel,
	}
	if err := initConfig(overrides); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting AionMCP server", 
		zap.String("version", "0.1.0"),
		zap.String("iteration", "0"))

	// Ensure data directory exists
	if err := ensureDataDirectory(); err != nil {
		logger.Fatal("Failed to create data directory", zap.Error(err))
	}

	// Create server instance
	server, err := core.NewServer(logger)
	if err != nil {
		logger.Fatal("Failed to create server", zap.Error(err))
	}

	// Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// Run server
	if err := server.Run(ctx); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}

	logger.Info("AionMCP server shutdown complete")
}

func initConfig(overrides ConfigOverrides) error {
	// Use custom config file if provided
	if overrides.ConfigFile != "" {
		viper.SetConfigFile(overrides.ConfigFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.grpc_port", 9090)
	viper.SetDefault("mcp.protocol_version", "1.0")
	viper.SetDefault("storage.type", "boltdb")
	viper.SetDefault("storage.path", "./data/aionmcp.db")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// Allow environment variable overrides
	viper.AutomaticEnv()
	viper.SetEnvPrefix("AIONMCP")

	// Override with command line flags if provided
	if overrides.HTTPPort != 0 {
		viper.Set("server.port", overrides.HTTPPort)
	}
	if overrides.GRPCPort != 0 {
		viper.Set("server.grpc_port", overrides.GRPCPort)
	}
	if overrides.LogLevel != "" {
		viper.Set("log.level", overrides.LogLevel)
	}

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

func initLogger() (*zap.Logger, error) {
	level := viper.GetString("log.level")
	format := viper.GetString("log.format")

	var config zap.Config
	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Parse log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}

func ensureDataDirectory() error {
	dataPath := viper.GetString("storage.path")
	if dataPath == "" {
		dataPath = "./data/aionmcp.db"
	}
	
	// Extract directory from path
	dir := filepath.Dir(dataPath)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory %s: %w", dir, err)
	}
	
	return nil
}