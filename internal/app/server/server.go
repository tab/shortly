package server

import (
	"net/http"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", HandleStatus)
	mux.HandleFunc("/", HandleCreateShortLink)
	mux.HandleFunc("/{id}", HandleGetShortLink)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
	return nil
}
