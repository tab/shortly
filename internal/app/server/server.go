package server

import "net/http"

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleCreateShortLink)
	mux.HandleFunc("/{id}", handleGetShortLink)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
	return nil
}
