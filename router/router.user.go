package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
)

func HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input api.UserParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "Failed to parse user"+err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := userService.Register(ctx, input)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "failed to register user", err)
		return
	}

	lib.Success(w, "user register successfully", userResponse)
}
