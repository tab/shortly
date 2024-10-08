package config

import (
	"flag"
	"os"
)

const (
	BaseURL       = "http://localhost:8080"
	ClientURL     = "http://localhost:3000"
	ServerAddress = "localhost:8080"
)

type AppConfig struct {
	Addr      string
	BaseURL   string
	ClientURL string
}

func New() *AppConfig {
	flagAddr := flag.String("a", ServerAddress, "address and port to run server")
	flagBaseURL := flag.String("b", BaseURL, "base address of the resulting shortened URL")
	flagClientURL := flag.String("c", ClientURL, "frontend client URL")
	flag.Parse()

	return &AppConfig{
		Addr:      getEnvOrFlag("SERVER_ADDRESS", *flagAddr),
		BaseURL:   getEnvOrFlag("BASE_URL", *flagBaseURL),
		ClientURL: getEnvOrFlag("CLIENT_URL", *flagClientURL),
	}
}

func getEnvOrFlag(envVar, flagValue string) string {
	if envValue, ok := os.LookupEnv(envVar); ok {
		return envValue
	}
	return flagValue
}
