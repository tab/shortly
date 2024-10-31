package config

import (
	"flag"
	"os"
)

const (
	BaseURL         = "http://localhost:8080"
	ClientURL       = "http://localhost:3000"
	ServerAddress   = "localhost:8080"
	FileStoragePath = "store.json"
)

type Config struct {
	Addr            string
	BaseURL         string
	ClientURL       string
	FileStoragePath string
}

func LoadConfig() *Config {
	flagAddr := flag.String("a", ServerAddress, "address and port to run server")
	flagBaseURL := flag.String("b", BaseURL, "base address of the resulting shortened URL")
	flagClientURL := flag.String("c", ClientURL, "frontend client URL")
	flagFileStoragePath := flag.String("f", FileStoragePath, "path to the file storage")
	flag.Parse()

	return &Config{
		Addr:            getEnvOrFlag("SERVER_ADDRESS", *flagAddr),
		BaseURL:         getEnvOrFlag("BASE_URL", *flagBaseURL),
		ClientURL:       getEnvOrFlag("CLIENT_URL", *flagClientURL),
		FileStoragePath: getEnvOrFlag("FILE_STORAGE_PATH", *flagFileStoragePath),
	}
}

func getEnvOrFlag(envVar, flagValue string) string {
	if envValue, ok := os.LookupEnv(envVar); ok {
		return envValue
	}
	return flagValue
}
