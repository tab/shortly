package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	endpoint := "http://localhost:8080"

	fmt.Println("Enter long URL:")

	reader := bufio.NewReader(os.Stdin)
	longUrl, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	longUrl = strings.TrimSuffix(longUrl, "\n")

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(longUrl).
		Post(endpoint)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status code: %s\n", response.Status())
	fmt.Printf("Short URL: %s\n", response.String())
}
