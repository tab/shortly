package errors

import "errors"

var (
	ErrorCouldNotGenerateCode    = errors.New("could not generate short code")
	ErrorInvalidURL              = errors.New("invalid URL")
	ErrorRequestBodyEmpty        = errors.New("request body is empty")
	ErrorShortLinkNotFound       = errors.New("short link not found")
	ErrorFailedToReadRandomBytes = errors.New("failed to read secure random bytes")
)
