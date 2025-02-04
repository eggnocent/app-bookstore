package helper

import (
	"app-bookstore/lib"
	"context"
	"net/http"
	"strings"
)

type KeyType string

const userKey KeyType = "user"

type Middleware struct {
	JWT lib.Jwt
}

func NewMiddleware(jwt lib.Jwt) *Middleware {
	return &Middleware{
		JWT: jwt,
	}
}

func (m *Middleware) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, lib.ErrMissingAuthHeader.Error(), http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			http.Error(w, lib.ErrMissingToken.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := m.JWT.VerifyAccessToken(token)
		if err != nil {
			http.Error(w, lib.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractClaims(r *http.Request) (*lib.JwtData, error) {
	claims, ok := r.Context().Value("user").(*lib.JwtData)
	if !ok || claims == nil {
		return nil, lib.ErrClaimsFound
	}
	return claims, nil
}
