package dto

import (
	"encoding/json"
	"io"
	"strings"

	"shortly/internal/app/errors"
	"shortly/internal/app/validator"
)

type CreateShortLinkParams struct {
	URL string `json:"url"`
}

type CreateShortLinkResponse struct {
	Result string `json:"result"`
}

type GetShortLinkResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (params *CreateShortLinkParams) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	return params.validateURL()
}

func (params *CreateShortLinkParams) DeprecatedValidate(body io.Reader) error {
	raw, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	params.URL = strings.Trim(strings.TrimSpace(string(raw)), "\"")

	return params.validateURL()
}

func (params *CreateShortLinkParams) validateURL() error {
	params.URL = strings.TrimSpace(params.URL)

	if params.URL == "" {
		return errors.ErrRequestBodyEmpty
	}

	if err := validator.Validate(params.URL); err != nil {
		return err
	}

	return nil
}
