package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"shortly/internal/logger"
)

// BaseURL is the base URL of the application
const BaseURL = "http://localhost:8080"

// ServerAddress is the address and port to run the server
const ServerAddress = "localhost:8080"

// ProfilerAddress is the address and port to run the profiler
const ProfilerAddress = "localhost:2080"

// EnableHTTPS is a flag to enable HTTPS
const EnableHTTPS = false

// Config is the application configuration
type Config struct {
	AppEnv          string
	Addr            string
	BaseURL         string
	ClientURL       string
	ProfilerAddr    string
	FileStoragePath string
	DatabaseDSN     string
	SecretKey       string
	EnableHTTPS     bool
	Certificate     string
	PrivateKey      string
}

// LoadConfig loads the application configuration
func LoadConfig() *Config {
	appLogger := logger.NewLogger()

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", env),
		fmt.Sprintf(".env.%s.local", env),
	}
	for _, file := range envFiles {
		err := godotenv.Overload(file)
		if err == nil {
			appLogger.Info().Msgf("Loaded %s file", file)
		}
	}

	flagAddr := flag.String("a", ServerAddress, "address and port to run server")
	flagBaseURL := flag.String("b", BaseURL, "base address of the resulting shortened URL")
	flagClientURL := flag.String("c", "", "frontend client URL")
	flagProfilerAddr := flag.String("p", ProfilerAddress, "address and port to run profiler")
	flagFileStoragePath := flag.String("f", "", "path to the file storage")
	flagDatabaseDSN := flag.String("d", "", "database DSN")
	flagSecretKey := flag.String("k", "", "JWT secret key")
	flagEnableHTTPS := flag.Bool("s", EnableHTTPS, "enable HTTPS")
	flag.Parse()

	return &Config{
		AppEnv:          env,
		Addr:            getEnvOrFlag("SERVER_ADDRESS", *flagAddr),
		BaseURL:         getEnvOrFlag("BASE_URL", *flagBaseURL),
		ClientURL:       getEnvOrFlag("CLIENT_URL", *flagClientURL),
		ProfilerAddr:    getEnvOrFlag("PROFILER_ADDRESS", *flagProfilerAddr),
		FileStoragePath: getEnvOrFlag("FILE_STORAGE_PATH", *flagFileStoragePath),
		DatabaseDSN:     getEnvOrFlag("DATABASE_DSN", *flagDatabaseDSN),
		SecretKey:       getEnvOrFlag("SECRET_KEY", *flagSecretKey),
		EnableHTTPS:     getEnvOrBoolFlag("ENABLE_HTTPS", *flagEnableHTTPS),
	}
}

func getEnvOrFlag(envVar, flagValue string) string {
	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}
	return flagValue
}

func getEnvOrBoolFlag(envVar string, flagValue bool) bool {
	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue == "true"
	}
	return flagValue
}
