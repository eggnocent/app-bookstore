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

		// cek token di header
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			lib.Error(w, http.StatusUnauthorized, "Missing or invalid authorization header", nil)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			lib.Error(w, http.StatusUnauthorized, "Invalid token", nil)
			return
		}

		// verif token
		claims, err := m.JWT.VerifyAccessToken(token)
		if err != nil {
			lib.Error(w, http.StatusUnauthorized, "Invalid or expired token", err)
			return
		}

		// cek token apa masih aktif bre
		sessionExists, err := model.CheckSessionExists(r.Context(), m.DB, token)
		if err != nil {
			lib.Error(w, http.StatusInternalServerError, "Failed to check session", err)
			return
		}

		if !sessionExists {
			lib.Error(w, http.StatusUnauthorized, "Session has expired, please login again", fmt.Errorf("session not found"))
			return
		}

		// ambil user_id dari token
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			lib.Error(w, http.StatusUnauthorized, "Invalid user ID", err)
			return
		}

		// ambil role_id dari user role
		roleID, err := model.GetUserRoleID(r.Context(), m.DB, userID)
		if err != nil {
			lib.Error(w, http.StatusInternalServerError, "Failed to retrieve user role ID", err)
			return
		}

		// cek role apakah memiliki akses endpoint
		requestedEndpoint := lib.NormalizeEndpoint(r.URL.Path)
		requestedMethod := r.Method

		allowed, err := model.CheckRoleAccess(r.Context(), m.DB, roleID, requestedEndpoint, requestedMethod)
		if err != nil {
			lib.Error(w, http.StatusInternalServerError, "Failed to check role access", err)
			return
		}

		if !allowed {
			lib.Error(w, http.StatusForbidden, "Access denied", fmt.Errorf("access denied for role %s on %s %s", roleID, requestedMethod, requestedEndpoint))
			return
		}

		// set ccontext untuk user
		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))

		// logging debug
		microSec := time.Since(start).Microseconds()
		fmt.Printf("[DEBUG] User %s accessed %s | Method: %s | Duration: %vμs\n", userID, requestedEndpoint, requestedMethod, microSec)
	})
}
