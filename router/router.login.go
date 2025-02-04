package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"
)

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input api.UserParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	loginResponse, err := userService.Login(ctx, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to login", err)
		return
	}

	lib.Success(w, "user login successfully", loginResponse)
}

func HandlerChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var input api.ChangePasswordParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = userService.ChangePassword(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "Failed to change password", err)
		return
	}

	lib.Success(w, "Password changed successfully", nil)
}
