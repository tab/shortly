package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"shortly/internal/app/dto"
	"shortly/internal/app/service"
)

// CookieName is the name of the authentication cookie
const CookieName = "auth"

// Middleware is a middleware for authentication
func Middleware(authenticator service.Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var currentUserID uuid.UUID
			cookie, err := r.Cookie(CookieName)

			if err != nil || cookie == nil {
				currentUserID, err = currentUser(w, authenticator)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				currentUserID, err = authenticator.Verify(cookie.Value)
				if err != nil {
					currentUserID, err = currentUser(w, authenticator)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}

			ctx := context.WithValue(r.Context(), dto.CurrentUser, currentUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func currentUser(w http.ResponseWriter, authenticator service.Authenticator) (uuid.UUID, error) {
	currentUserID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	token, err := authenticator.Generate(currentUserID)
	if err != nil {
		return uuid.Nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	return currentUserID, nil
}
