package errors

import "errors"

var (
	// ErrInvalidURL is returned when the URL is invalid
	ErrInvalidURL = errors.New("invalid URL")

	// ErrOriginalURLEmpty is returned when the original URL is empty
	ErrOriginalURLEmpty = errors.New("original URL is required")

	// ErrCorrelationIDEmpty is returned when the correlation ID is empty
	ErrCorrelationIDEmpty = errors.New("correlation id is required")

	// ErrShortCodeEmpty is returned when the short code is empty
	ErrShortCodeEmpty = errors.New("short code is required")

	// ErrShortLinkNotFound is returned when the short link is not found
	ErrShortLinkNotFound = errors.New("short link not found")

	// ErrShortLinkDeleted is returned when the short link is deleted
	ErrShortLinkDeleted = errors.New("short link deleted")

	// ErrFailedToReadRandomBytes is returned when the random bytes cannot be read
	ErrFailedToReadRandomBytes = errors.New("failed to read secure random bytes")

	// ErrFailedToGenerateCode is returned when the short code cannot be generated
	ErrFailedToGenerateCode = errors.New("failed to generate short code")

	// ErrFailedToGenerateUUID is returned when the UUID cannot be generated
	ErrFailedToGenerateUUID = errors.New("failed to generate UUID")

	// ErrFailedToOpenFile is returned when the file cannot be opened
	ErrFailedToOpenFile = errors.New("failed to open file")

	// ErrorFailedToReadFromFile is returned when the file cannot be read
	ErrorFailedToReadFromFile = errors.New("failed to read from file")

	// ErrFailedToWriteToFile is returned when the file cannot be written to
	ErrFailedToWriteToFile = errors.New("failed to write to file")

	// ErrFailedToSaveURL is returned when the URL cannot be saved
	ErrFailedToSaveURL = errors.New("failed to save URL")

	// ErrURLAlreadyExists is returned when the URL already exists
	ErrURLAlreadyExists = errors.New("URL already exists")

	// ErrInvalidToken is returned when JWT token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidSigningMethod is returned when the signing method is invalid
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	// ErrInvalidUserID is returned when the user ID is invalid
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrFailedToLoadUserUrls is returned when the user URLs cannot be loaded
	ErrFailedToLoadUserUrls = errors.New("failed to load user URLs")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
