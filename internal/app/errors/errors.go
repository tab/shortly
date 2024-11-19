package errors

import "errors"

var (
	ErrInvalidURL              = errors.New("invalid URL")
	ErrOriginalURLEmpty        = errors.New("original URL is required")
	ErrCorrelationIDEmpty      = errors.New("correlation id is required")
	ErrShortLinkNotFound       = errors.New("short link not found")
	ErrFailedToReadRandomBytes = errors.New("failed to read secure random bytes")
	ErrFailedToGenerateCode    = errors.New("failed to generate short code")
	ErrFailedToGenerateUUID    = errors.New("failed to generate UUID")
	ErrFailedToOpenFile        = errors.New("failed to open file")
	ErrorFailedToReadFromFile  = errors.New("failed to read from file")
	ErrFailedToWriteToFile     = errors.New("failed to write to file")
	ErrFailedToSaveURL         = errors.New("failed to save URL")
	ErrURLAlreadyExists        = errors.New("URL already exists")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidSigningMethod    = errors.New("invalid signing method")
	ErrInvalidUserID           = errors.New("invalid user id")
	ErrFailedToLoadUserUrls    = errors.New("failed to load user URLs")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
