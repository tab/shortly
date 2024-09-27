package server

import (
	"fmt"
	"io"
	"net/http"
	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
	"strings"
)

var storage = store.NewURLStore()

func HandleCreateShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
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

	longURL := string(body)
	shortCode := helpers.GenerateShortCode()
	shortURL := fmt.Sprintf(`http://localhost:8080/%s`, shortCode)

	storage.Set(shortCode, longURL)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func HandleGetShortLink(res http.ResponseWriter, req *http.Request) {
	shortCode := strings.TrimPrefix(req.URL.Path, "/")

	longURL, found := storage.Get(shortCode)
	if !found {
		http.Error(res, "Short code not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	http.Redirect(res, req, longURL, http.StatusTemporaryRedirect)
}
