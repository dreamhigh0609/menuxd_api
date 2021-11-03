package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
)

// Errors
var (
	ErrFailedSigning          = errors.New("The token could not be signed")
	ErrTokenExpired           = errors.New("Your token has expired")
	ErrSignatureNoMatch       = errors.New("Token signature does not match")
	ErrTokenNoValid           = errors.New("Your token is not valid")
	ErrUserNotAuthorized      = errors.New("Not authorized")
	ErrInsufficientPrivileges = errors.New("Insufficient privileges")
	ErrTokenNotFound          = errors.New("Token is required")
)

func Authenticator(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())

			if err != nil {
				http.Error(
					w,
					ErrTokenNotFound.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			if token == nil || !token.Valid {
				http.Error(
					w,
					ErrTokenNoValid.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			claimRole, ok := claims["role"].(string)
			if !ok {
				http.Error(
					w,
					ErrInsufficientPrivileges.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			if role != claimRole && claimRole != "admin" {
				http.Error(
					w,
					ErrInsufficientPrivileges.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}
