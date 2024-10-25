package dto

import (
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

func (r *CreateShortLinkParams) Validate() error {
  r.URL = strings.TrimSpace(r.URL)

  if r.URL == "" {
    return errors.ErrRequestBodyEmpty
  }

  if err := validator.Validate(r.URL); err != nil {
    return err
  }

  return nil
}
