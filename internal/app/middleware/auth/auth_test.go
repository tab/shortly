package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/app/service"
)

func Test_Auth_Middleware(t *testing.T) {
	client := &http.Client{}
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	authenticator := service.NewAuthService(cfg)

	type result struct {
		code            int
		cookieGenerated bool
	}

	tests := []struct {
		name        string
		isProtected bool
		cookie      *http.Cookie
		expected    result
	}{
		{
			name:        "Public route, no cookie",
			isProtected: false,
			cookie:      nil,
			expected: result{
				code:            http.StatusOK,
				cookieGenerated: true,
			},
		},
		{
			name:        "Public route, valid cookie",
			isProtected: false,
			cookie: func() *http.Cookie {
				c, _ := generateCookie(authenticator)
				return c
			}(),
			expected: result{
				code:            http.StatusOK,
				cookieGenerated: false,
			},
		},
		{
			name:        "Public route, invalid cookie",
			isProtected: false,
			cookie:      &http.Cookie{Name: CookieName, Value: "invalid-token"},
			expected: result{
				code:            http.StatusOK,
				cookieGenerated: true,
			},
		},
		{
			name:        "Protected route, no cookie",
			isProtected: true,
			cookie:      nil,
			expected: result{
				code:            http.StatusOK,
				cookieGenerated: true,
			},
		},
		{
			name:        "Protected route, valid cookie",
			isProtected: true,
			cookie: func() *http.Cookie {
				c, _ := generateCookie(authenticator)
				return c
			}(),
			expected: result{
				code:            http.StatusOK,
				cookieGenerated: false,
			},
		},
		{
			name:        "Protected route, invalid cookie",
			isProtected: true,
			cookie:      &http.Cookie{Name: CookieName, Value: "invalid-token"},
			expected: result{
				code:            http.StatusOK,
				cookieGenerated: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"response": "ok"}`))
				assert.NoError(t, err)
			})

			var wrappedHandler http.Handler = handler
			wrappedHandler = Middleware(authenticator)(wrappedHandler)

			ts := httptest.NewServer(wrappedHandler)
			defer ts.Close()

			req, err := http.NewRequest("GET", ts.URL, nil)
			assert.NoError(t, err)

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expected.code, resp.StatusCode)

			cookieFound := false
			for _, c := range resp.Cookies() {
				if c.Name == CookieName {
					cookieFound = true
					break
				}
			}
			assert.Equal(t, tt.expected.cookieGenerated, cookieFound)
		})
	}
}

func generateCookie(authenticator service.Authenticator) (*http.Cookie, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	token, err := authenticator.Generate(id)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}, nil
}
