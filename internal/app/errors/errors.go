package errors

import "errors"

var (
	ErrInvalidURL                   = errors.New("invalid URL")
	ErrRequestBodyEmpty             = errors.New("request body is empty")
	ErrShortLinkNotFound            = errors.New("short link not found")
	ErrFailedToReadRandomBytes      = errors.New("failed to read secure random bytes")
	ErrFailedToGenerateCode         = errors.New("failed to generate short code")
	ErrFailedToGenerateUUID         = errors.New("failed to generate UUID")
	ErrFilePathIsDirectory          = errors.New("file path is a directory")
	ErrFailedToOpenFile             = errors.New("failed to open file")
	ErrFailedToWriteToFile          = errors.New("failed to write to file")
	ErrFailedToInitializeRepository = errors.New("failed to initialize repository")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
