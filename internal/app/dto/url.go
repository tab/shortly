package dto

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/google/uuid"

	"shortly/internal/app/errors"
	"shortly/internal/app/validator"
)

// CreateShortLinkRequest is a request for short link creation
type CreateShortLinkRequest struct {
	URL string `json:"url"`
}

// CreateShortLinkResponse is a response for short link creation
type CreateShortLinkResponse struct {
	Result string `json:"result"`
}

// BatchCreateShortLinkParams is a request for batch short link creation
type BatchCreateShortLinkParams struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchCreateShortLinkRequest is a request for batch short link creation
type BatchCreateShortLinkRequest []BatchCreateShortLinkParams

// BatchCreateShortLinkResponse is a response for batch short link creation
type BatchCreateShortLinkResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// BatchCreateShortLinkResponses is a response for batch short link creation
type BatchCreateShortLinkResponses []BatchCreateShortLinkResponse

// GetShortLinkRequest is a request for short link retrieval
type GetShortLinkResponse struct {
	Result string `json:"result"`
}

// GetUserURLsResponse is a response for user URLs retrieval
type GetUserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// BatchDeleteShortLinkRequest is a request for batch short link deletion
type BatchDeleteShortLinkRequest []string

// BatchDeleteParams is a request for batch short link deletion
type BatchDeleteParams struct {
	UserID     uuid.UUID
	ShortCodes []string
}

// ErrorResponse is a response for batch short link deletion
type ErrorResponse struct {
	Error string `json:"error"`
}

// Validate validates a create short link request
func (params *CreateShortLinkRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	return params.validateURL()
}

// Validate validates a batch create short link request
func (params *BatchCreateShortLinkRequest) Validate(body io.Reader) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(params); err != nil {
		return err
	}

	for _, p := range *params {
		p.CorrelationID = strings.TrimSpace(p.CorrelationID)
		if p.CorrelationID == "" {
			return errors.ErrCorrelationIDEmpty
		}

		if err := validator.Validate(p.OriginalURL); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates a batch delete short link request
func (params *BatchDeleteShortLinkRequest) Validate(body io.Reader) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(params); err != nil {
		return err
	}

	if len(*params) == 0 {
		return errors.ErrShortCodeEmpty
	}

	for _, p := range *params {
		if p == "" {
			return errors.ErrShortCodeEmpty
		}
	}

	return nil
}

// DeprecatedValidate validates a create short link request (text/plain endpoint)
func (params *CreateShortLinkRequest) DeprecatedValidate(body io.Reader) error {
	raw, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	params.URL = strings.Trim(strings.TrimSpace(string(raw)), "\"")

	return params.validateURL()
}

// validateURL validates a URL
func (params *CreateShortLinkRequest) validateURL() error {
	params.URL = strings.TrimSpace(params.URL)

	if params.URL == "" {
		return errors.ErrOriginalURLEmpty
	}

	if err := validator.Validate(params.URL); err != nil {
		return err
	}

	return nil
}
