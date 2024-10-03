package config

import (
	"flag"
	"os"
)

const ServerAddress = "localhost:8080"
const BaseURL = "http://localhost:8080"

type Options struct {
	Addr    string
	BaseURL string
}

var options = Options{
	Addr:    ServerAddress,
	BaseURL: BaseURL,
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

func Init() Options {
	flagAddr := flag.String("a", options.Addr, "address and port to run server")
	flagBaseURL := flag.String("b", options.BaseURL, "base address of the resulting shortened URL")
	flag.Parse()

	options.Addr = setServerAddress(*flagAddr)
	options.BaseURL = setBaseURL(*flagBaseURL)

	return options
}
