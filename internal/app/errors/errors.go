package errors

import "errors"

var (
	ErrCouldNotGenerateCode    = errors.New("could not generate short code")
	ErrInvalidURL              = errors.New("invalid URL")
	ErrRequestBodyEmpty        = errors.New("request body is empty")
	ErrShortLinkNotFound       = errors.New("short link not found")
	ErrFailedToReadRandomBytes = errors.New("failed to read secure random bytes")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
