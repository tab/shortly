package service

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
)

func TestCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockURLRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	service := NewURLService(repo, rand, cfg)

	type result struct {
		shortCode string
		shortURL  string
		error     error
	}

	tests := []struct {
		name     string
		body     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			body: "https://example.com",
			before: func() {
				rand.EXPECT().Hex().Return("abcd1234", nil)

				url := repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}
				repo.EXPECT().Set(url)
			},
			expected: result{
				shortCode: "abcd1234",
				shortURL:  "http://localhost:8080/abcd1234",
				error:     nil,
			},
		},
		{
			name:   "Empty Body",
			body:   "",
			before: func() {},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrorRequestBodyEmpty,
			},
		},
		{
			name:   "Invalid URL",
			body:   "not-a-url",
			before: func() {},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrorInvalidURL,
			},
		},
		{
			name: "Error generating short code",
			body: "https://example.com",
			before: func() {
				rand.EXPECT().Hex().Return("", errors.ErrorFailedToReadRandomBytes)
			},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrorCouldNotGenerateCode,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.before()

			req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			result, err := service.CreateShortLink(req)

			assert.Equal(t, test.expected.shortURL, result)
			assert.Equal(t, test.expected.error, err)
		})
	}

}

func TestGetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockURLRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	service := NewURLService(repo, rand, cfg)

	type result struct {
		url   *repository.URL
		found bool
	}

	tests := []struct {
		name      string
		shortCode string
		expected  result
	}{
		{
			name:      "Success",
			shortCode: "abcd1234",
			expected: result{
				url: &repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				},
				found: true,
			},
		},
		{
			name:      "Not Found",
			shortCode: "1234abcd",
			expected: result{
				url:   nil,
				found: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo.EXPECT().Get(test.shortCode).Return(test.expected.url, test.expected.found)

			url, found := service.GetShortLink(test.shortCode)

			assert.Equal(t, test.expected.url, url)
			assert.Equal(t, test.expected.found, found)
		})
	}
}
