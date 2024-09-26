package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	. "shortly/internal/app/helpers"
)

func handleCreateShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		fmt.Println("Error reading body:", err)
		return
	}
	defer req.Body.Close()

	longURL := string(body)
	fmt.Println("Received URL to shorten:", longURL)

	shortID := GenerateShortID()
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

	fmt.Println("Generated short URL:", shortURL)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func handleGetShortLink(res http.ResponseWriter, req *http.Request) {
}
