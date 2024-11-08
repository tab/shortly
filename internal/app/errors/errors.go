package errors

import "errors"

var (
	ErrInvalidURL              = errors.New("invalid URL")
	ErrRequestBodyEmpty        = errors.New("request body is empty")
	ErrShortLinkNotFound       = errors.New("short link not found")
	ErrFailedToReadRandomBytes = errors.New("failed to read secure random bytes")
	ErrFailedToGenerateCode    = errors.New("failed to generate short code")
	ErrFailedToGenerateUUID    = errors.New("failed to generate UUID")
	ErrFailedToOpenFile        = errors.New("failed to open file")
	ErrorFailedToReadFromFile  = errors.New("failed to read from file")
	ErrFailedToWriteToFile     = errors.New("failed to write to file")
	ErrFailedToSaveURL         = errors.New("failed to save URL")
	ErrFileStoragePathNotSet   = errors.New("file storage path not set")
	ErrMementoNotImplemented   = errors.New("memento not implemented")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
