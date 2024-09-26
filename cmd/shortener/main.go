package main

import "shortly/internal/app/server"

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
