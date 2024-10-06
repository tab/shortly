package helpers

import "regexp"

const URLRegex = `^(https?:\/\/)?([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,6}([\/\w .-]*)*\/?$`

func IsValidURL(url string) bool {
	return regexp.MustCompile(URLRegex).MatchString(url)
}

func IsInvalidURL(url string) bool {
	return !IsValidURL(url)
}
