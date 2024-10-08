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
		Addr:      setServerAddress(*flagAddr),
		BaseURL:   setBaseURL(*flagBaseURL),
		ClientURL: setClientURL(*flagClientURL),
	}
}

func setServerAddress(flagAddr string) string {
	envAddr, ok := os.LookupEnv("SERVER_ADDRESS")

	if ok {
		return envAddr
	}

	return flagAddr
}

func setBaseURL(flagBaseURL string) string {
	envBaseURL, ok := os.LookupEnv("BASE_URL")

	if ok {
		return envBaseURL
	}

	return flagBaseURL
}

func setClientURL(flagClientURL string) string {
	envClientURL, ok := os.LookupEnv("CLIENT_URL")

	if ok {
		return envClientURL
	}

	return flagClientURL
}
