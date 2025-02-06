package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerUserRequestList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRequest, err := userRequestService.ListUserRequest(ctx)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retreive user request", err)
		return
	}

	lib.Success(w, "successfully to retreive user request", userRequest)
}

func HandlerUserRequestDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user request id", http.StatusBadRequest)
		return
	}

	userRequestResponse, err := userRequestService.GetOneUserRequest(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retreive user request detail", err)
		return
	}

	lib.Success(w, "user request detail successfully retreived", userRequestResponse)
}

func HandlerUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	var input api.UserRequestParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userResponse, err := userRequestService.CreateUserRequest(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to request role to admin", err)
		return
	}

	lib.Success(w, "success to request role to admin", userResponse)
}

func HandlerUserRequestUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "invalid or expired token", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "invalid token", http.StatusInternalServerError)
		return
	}

	var input api.UserRequestUpdateParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body"+err.Error(), http.StatusInternalServerError)
		return
	}

	userReqUpdate, err := userRequestService.UpdateUserRequest(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update user request", err)
		return
	}

	lib.Success(w, "role successfully updated", userReqUpdate)
}
