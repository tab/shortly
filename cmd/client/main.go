package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	BaseURL = "http://localhost:8080"
)

type Config struct {
	BaseURL string
}

func main() {
	cfg := parseFlags()
	longURL, err := readLongURL()
	if err != nil {
		log.Fatalf("Error reading URL: %v", err)
	}

	shortURL, err := createShortLink(cfg.BaseURL, longURL)
	if err != nil {
		log.Fatalf("Error creating short link: %v", err)
	}

	fmt.Printf("Short URL: %s\n", shortURL)
}

func parseFlags() Config {
	var cfg Config
	flag.StringVar(&cfg.BaseURL, "a", BaseURL, "address and port to call server")
	flag.Parse()
	return cfg
}

func readLongURL() (string, error) {
	fmt.Println("Enter long URL:")
	reader := bufio.NewReader(os.Stdin)
	longURL, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(longURL, "\n"), nil
}

func createShortLink(endpoint, longURL string) (string, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(longURL).
		Post(endpoint)

	if err != nil {
		return "", err
	}

	if response.IsError() {
		return "", fmt.Errorf("server responded with status code %d", response.StatusCode())
	}

	return response.String(), nil
}
