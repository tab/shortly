package errors

import "errors"

var (
	ErrInvalidURL              = errors.New("invalid URL")
	ErrRequestBodyEmpty        = errors.New("request body is empty")
	ErrShortLinkNotFound       = errors.New("short link not found")
	ErrFailedToReadRandomBytes = errors.New("failed to read secure random bytes")
	ErrFailedToGenerateCode    = errors.New("failed to generate short code")
	ErrFailedToGenerateUUID    = errors.New("failed to generate UUID")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
