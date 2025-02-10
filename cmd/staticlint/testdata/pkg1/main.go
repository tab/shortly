package main

import "os"

func callExit() {
	os.Exit(4) // want "direct call to os.Exit found in ..."
}

func main() {
	os.Exit(1) // want "direct call to os.Exit found in ..."

	defer os.Exit(2) // want "direct call to os.Exit found in ..."

	if true {
		os.Exit(3) // want "direct call to os.Exit found in ..."
	}

	callExit()
}
