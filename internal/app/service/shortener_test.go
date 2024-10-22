package service

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
)

func Test_CreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockURLRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	service := NewURLService(cfg, repo, rand)

	type result struct {
		shortCode string
		shortURL  string
		error     error
	}

	tests := []struct {
		name     string
		body     io.Reader
		before   func()
		expected result
	}{
		{
			name: "Success",
			body: strings.NewReader(`{"url":"https://example.com"}`),
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
			body:   strings.NewReader(""),
			before: func() {},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrRequestBodyEmpty,
			},
		},
		{
			name:   "Invalid URL",
			body:   strings.NewReader(`{"url":"not-a-url"}`),
			before: func() {},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrInvalidURL,
			},
		},
		{
			name: "Error generating short code",
			body: strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrFailedToGenerateCode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req, _ := http.NewRequest(http.MethodPost, "/", tt.body)
			shortURL, err := service.CreateShortLink(req)

			assert.Equal(t, tt.expected.shortURL, shortURL)
			assert.Equal(t, tt.expected.error, err)
		})
	}
}

func Test_DeprecatedCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockURLRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	service := NewURLService(cfg, repo, rand)

	type result struct {
		shortCode string
		shortURL  string
		error     error
	}

	tests := []struct {
		name     string
		body     io.Reader
		before   func()
		expected result
	}{
		{
			name: "Success",
			body: strings.NewReader("https://example.com"),
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
			body:   strings.NewReader(""),
			before: func() {},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrRequestBodyEmpty,
			},
		},
		{
			name:   "Invalid URL",
			body:   strings.NewReader("not-a-url"),
			before: func() {},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrInvalidURL,
			},
		},
		{
			name: "Error generating short code",
			body: strings.NewReader("https://example.com"),
			before: func() {
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrFailedToGenerateCode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req, _ := http.NewRequest(http.MethodPost, "/", tt.body)
			shortURL, err := service.DeprecatedCreateShortLink(req)

			assert.Equal(t, tt.expected.shortURL, shortURL)
			assert.Equal(t, tt.expected.error, err)
		})
	}

}

func Test_GetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockURLRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	service := NewURLService(cfg, repo, rand)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().Get(tt.shortCode).Return(tt.expected.url, tt.expected.found)

			url, found := service.GetShortLink(tt.shortCode)

			assert.Equal(t, tt.expected.url, url)
			assert.Equal(t, tt.expected.found, found)
		})
	}
}
