package errors

import "fmt"

type MethodNotAllowedError struct{}

func (e *MethodNotAllowedError) Error() string {
	return "Wrong HTTP method"
}

type ResponseWriteError struct{}

func (e *ResponseWriteError) Error() string {
	return "Failed to write response"
}

type InvalidRequestBodyError struct{}

func (e *InvalidRequestBodyError) Error() string {
	return "Unable to process request"
}

type InvalidURLError struct {
	URL string
}

func (e *InvalidURLError) Error() string {
	return fmt.Sprintf("Invalid URL: %s", e.URL)
}

type ShortCodeGenerationError struct{}

func (e *ShortCodeGenerationError) Error() string {
	return "Failed to generate short code"
}

type ShortCodeNotFoundError struct {
	Code string
}

func (e *ShortCodeNotFoundError) Error() string {
	return fmt.Sprintf("Short code not found: %s", e.Code)
}
