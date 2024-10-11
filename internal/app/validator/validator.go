package validator

import (
	"regexp"

	"shortly/internal/app/errors"
)

const URLRegex = `^(https?:\/\/)([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,6}([\/\w .-]*)*\/?$`

var urlRegex = regexp.MustCompile(URLRegex)

func Validate(url string) error {
	if !urlRegex.MatchString(url) {
		return errors.ErrInvalidURL
	}
	return nil
}
