package config

import (
	"flag"
	"os"
)

const ServerAddress = "localhost:8080"
const BaseURL = "http://localhost:8080"
const ClientURL = "http://localhost:3000"

type Options struct {
	Addr      string
	BaseURL   string
	ClientURL string
}

var options = Options{
	Addr:      ServerAddress,
	BaseURL:   BaseURL,
	ClientURL: ClientURL,
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

func Init() Options {
	flagAddr := flag.String("a", options.Addr, "address and port to run server")
	flagBaseURL := flag.String("b", options.BaseURL, "base address of the resulting shortened URL")
	flagClientURL := flag.String("c", options.ClientURL, "frontend client URL")
	flag.Parse()

	options.Addr = setServerAddress(*flagAddr)
	options.BaseURL = setBaseURL(*flagBaseURL)
	options.ClientURL = setClientURL(*flagClientURL)

	return options
}
