package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"shortly/internal/app/config"
)

var options config.Options

func AppRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.Heartbeat("/status"),
		middleware.Logger,
		middleware.RequestID,
		middleware.Recoverer)

	router.Post("/", HandleCreateShortLink)
	router.Get("/{id}", HandleGetShortLink)

	return router
}

func Run() {
	options = config.Init()

	fmt.Println("Running server on", options.Addr)
	fmt.Println("Shortener base address is", options.BaseURL)

	err := http.ListenAndServe(options.Addr, AppRouter())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
