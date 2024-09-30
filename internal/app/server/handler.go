package server

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
)

var storage = store.NewURLStore()
var urlPattern = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

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
	if len(body) == 0 {
		http.Error(res, "Unable to process request", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	longURL := string(body)
	if !urlPattern.MatchString(longURL) {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

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
