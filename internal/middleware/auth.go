package middleware

import (
	"github.com/vindosVP/loyalty-system/pkg/auth"
	"github.com/vindosVP/loyalty-system/pkg/tokens"
	"net/http"
)

type Authenticator struct {
	secret string
}

func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{secret: secret}
}

func (a *Authenticator) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := auth.ParseBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		authorized, err := tokens.IsAuthorized(token, a.secret)
		if err != nil || !authorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := tokens.ExtractID(token, a.secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r.Header.Set("x-user-id", userID)
		next.ServeHTTP(w, r)
	})
}
