package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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
	log.Fatal(http.ListenAndServe(":8080", AppRouter()))
}
