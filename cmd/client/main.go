package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Options struct {
	Endpoint string
}

var opts Options

func main() {
	flag.StringVar(&opts.Endpoint, "a", "http://localhost:8080", "address and port to call server")
	flag.Parse()

	fmt.Println("Enter long URL:")

	reader := bufio.NewReader(os.Stdin)
	longURL, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	longURL = strings.TrimSuffix(longURL, "\n")

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(longURL).
		Post(opts.Endpoint)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status code: %s\n", response.Status())
	fmt.Printf("Short URL: %s\n", response.String())
}
