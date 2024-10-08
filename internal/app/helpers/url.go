package helpers

import (
	"regexp"

	"shortly/internal/app/errors"
)

const URLRegex = `^(https?:\/\/)?([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,6}([\/\w .-]*)*\/?$`

var urlRegex = regexp.MustCompile(URLRegex)

func Validate(url string) (bool, error) {
	isValid := urlRegex.MatchString(url)
	if !isValid {
		return isValid, &errors.InvalidURLError{URL: url}
	}
	return isValid, nil
}
