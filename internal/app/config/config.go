package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"shortly/internal/logger"
)

const (
	BaseURL       = "http://localhost:8080"
	ServerAddress = "localhost:8080"
)

// Config is the application configuration
type Config struct {
	AppEnv          string
	Addr            string
	BaseURL         string
	ClientURL       string
	FileStoragePath string
	DatabaseDSN     string
	SecretKey       string
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
	flagFileStoragePath := flag.String("f", "", "path to the file storage")
	flagDatabaseDSN := flag.String("d", "", "database DSN")
	flagSecretKey := flag.String("s", "", "JWT secret key")
	flag.Parse()

	return &Config{
		AppEnv:          env,
		Addr:            getEnvOrFlag("SERVER_ADDRESS", *flagAddr),
		BaseURL:         getEnvOrFlag("BASE_URL", *flagBaseURL),
		ClientURL:       getEnvOrFlag("CLIENT_URL", *flagClientURL),
		FileStoragePath: getEnvOrFlag("FILE_STORAGE_PATH", *flagFileStoragePath),
		DatabaseDSN:     getEnvOrFlag("DATABASE_DSN", *flagDatabaseDSN),
		SecretKey:       getEnvOrFlag("SECRET_KEY", *flagSecretKey),
	}
}

func getEnvOrFlag(envVar, flagValue string) string {
	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}
	return flagValue
}
