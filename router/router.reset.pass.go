package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
)

func HandlerRequestResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input api.PasswordResetRequestParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := passwordResetService.RequestReset(ctx, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "Failed to request reset passworrd", err)
		return
	}

	lib.Success(w, "password reset request", response)
}

func HandlerValidateResetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")

	response, err := passwordResetService.ValidateToken(ctx, token)
	if err != nil {
		lib.Error(w, http.StatusUnauthorized, "invalid or expired token", err)
		return
	}

	lib.Success(w, "token is valid", response)
}

func HandlerResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input api.PasswordResetConfigParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := passwordResetService.ResetPassword(ctx, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to reset password", err)
		return
	}

	lib.Success(w, "password successfully reset", response)
}
