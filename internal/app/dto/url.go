package dto

import (
	"encoding/json"
	"io"
	"strings"

	"shortly/internal/app/errors"
	"shortly/internal/app/validator"
)

type CreateShortLinkRequest struct {
	URL string `json:"url"`
}

type CreateShortLinkResponse struct {
	Result string `json:"result"`
}

type BatchCreateShortLinkParams struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchCreateShortLinkRequest []BatchCreateShortLinkParams

type BatchCreateShortLinkResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchCreateShortLinkResponses []BatchCreateShortLinkResponse

type GetShortLinkResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (params *CreateShortLinkRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	return params.validateURL()
}

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

func (params *CreateShortLinkRequest) DeprecatedValidate(body io.Reader) error {
	raw, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	params.URL = strings.Trim(strings.TrimSpace(string(raw)), "\"")

	return params.validateURL()
}

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
