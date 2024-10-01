package config

import (
	"flag"
)

type Options struct {
	Addr    string
	BaseURL string
}

var options Options

func Init() Options {
	flag.StringVar(&options.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&options.BaseURL, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.Parse()

	return options
}
