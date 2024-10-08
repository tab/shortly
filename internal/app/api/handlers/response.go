package handlers

import (
	"net/http"

	"shortly/internal/app/errors"
)

func httpResponse(res http.ResponseWriter, statusCode int, body []byte, redirectURL string) {
	res.Header().Set("Content-Type", "text/plain")

	if redirectURL != "" {
		http.Redirect(res, nil, redirectURL, statusCode)
		return
	}

	res.WriteHeader(statusCode)
	if body != nil {
		_, err := res.Write(body)
		if err != nil {
			httpError(res, &errors.ResponseWriteError{}, http.StatusInternalServerError)
		}
	}
}

func httpError(res http.ResponseWriter, err error, statusCode int) {
	res.Header().Set("Content-Type", "text/plain")

	http.Error(res, err.Error(), statusCode)
}
