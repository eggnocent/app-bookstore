package helper

import (
	"app-bookstore/lib"
	"app-bookstore/model"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type KeyType string

const userKey KeyType = "user"

type Middleware struct {
	JWT lib.Jwt
	DB  *sqlx.DB
}

func NewMiddleware(jwt lib.Jwt, db *sqlx.DB) *Middleware {
	return &Middleware{
		JWT: jwt,
		DB:  db,
	}
}

func (m *Middleware) CheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		authHeader := r.Header.Get("Authorization")

		// cek token
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

		// verif token jwt
		claims, err := m.JWT.VerifyAccessToken(token)
		if err != nil {
			http.Error(w, lib.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		// ambil user id dari token
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusUnauthorized)
			return
		}

		// ambil role id dari user_roles
		roleID, err := model.GetUserRoleID(r.Context(), m.DB, userID)
		if err != nil {
			http.Error(w, "failed to retrieve user role ID", http.StatusInternalServerError)
			return
		}

		// check akses endpoint
		requestedEndpoint := lib.NormalizeEndpoint(r.URL.Path)
		requestedMethod := r.Method

		allowed, err := model.CheckRoleAccess(r.Context(), m.DB, roleID, requestedEndpoint, requestedMethod)

		if err != nil || !allowed {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))

		microSec := time.Since(start).Microseconds()
		fmt.Printf("User %s accessed %s %s | Duration: %vμs\n", r.URL, microSec)
	})
}
