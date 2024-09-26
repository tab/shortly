package server

import (
	"fmt"
	"io"
	"net/http"
	"shortly/internal/app/helpers"
)

func HandleCreateShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(res, "Unable to close reader", http.StatusInternalServerError)
		}
	}(req.Body)

	shortID := helpers.GenerateShortID()
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func HandleGetShortLink(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	http.Redirect(res, req, "https://practicum.yandex.ru/", http.StatusTemporaryRedirect)
}
