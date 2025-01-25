package validator

import (
	"net/url"

	"shortly/internal/app/errors"
)

// Validate checks if the URL is valid
func Validate(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(rawURL)

	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.ErrInvalidURL
	}

	return nil
}
