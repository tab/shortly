package errors

import "errors"

// ErrInvalidURL is returned when the URL is invalid
var ErrInvalidURL = errors.New("invalid URL")

// ErrOriginalURLEmpty is returned when the original URL is empty
var ErrOriginalURLEmpty = errors.New("original URL is required")

// ErrCorrelationIDEmpty is returned when the correlation ID is empty
var ErrCorrelationIDEmpty = errors.New("correlation id is required")

// ErrInvalidBatchParams is returned when the batch params are invalid
var ErrInvalidBatchParams = errors.New("invalid batch params")

// ErrInvalidShortCode is returned when the short code is invalid
var ErrInvalidShortCode = errors.New("invalid short code")

// ErrShortCodeEmpty is returned when the short code is empty
var ErrShortCodeEmpty = errors.New("short code is required")

// ErrShortLinkNotFound is returned when the short link is not found
var ErrShortLinkNotFound = errors.New("short link not found")

// ErrShortLinkDeleted is returned when the short link is deleted
var ErrShortLinkDeleted = errors.New("short link deleted")

// ErrFailedToReadRandomBytes is returned when the random bytes cannot be read
var ErrFailedToReadRandomBytes = errors.New("failed to read secure random bytes")

// ErrFailedToGenerateCode is returned when the short code cannot be generated
var ErrFailedToGenerateCode = errors.New("failed to generate short code")

// ErrFailedToGenerateUUID is returned when the UUID cannot be generated
var ErrFailedToGenerateUUID = errors.New("failed to generate UUID")

// ErrFailedToOpenFile is returned when the file cannot be opened
var ErrFailedToOpenFile = errors.New("failed to open file")

// ErrorFailedToReadFromFile is returned when the file cannot be read
var ErrorFailedToReadFromFile = errors.New("failed to read from file")

// ErrFailedToWriteToFile is returned when the file cannot be written to
var ErrFailedToWriteToFile = errors.New("failed to write to file")

// ErrFailedToSaveURL is returned when the URL cannot be saved
var ErrFailedToSaveURL = errors.New("failed to save URL")

// ErrURLAlreadyExists is returned when the URL already exists
var ErrURLAlreadyExists = errors.New("URL already exists")

// ErrInvalidToken is returned when JWT token is invalid
var ErrInvalidToken = errors.New("invalid token")

// ErrInvalidSigningMethod is returned when the signing method is invalid
var ErrInvalidSigningMethod = errors.New("invalid signing method")

// ErrInvalidUserID is returned when the user ID is invalid
var ErrInvalidUserID = errors.New("invalid user id")

// ErrFailedToLoadUserUrls is returned when the user URLs cannot be loaded
var ErrFailedToLoadUserUrls = errors.New("failed to load user URLs")

// Is a shortcut for errors.Is
var Is = errors.Is
