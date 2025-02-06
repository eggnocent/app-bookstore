package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRoleResponse, err := userRolesService.ListUserRole(ctx)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve user roles", err)
		return
	}
	lib.Success(w, "success to retreive user roles", userRoleResponse)
}

func HandlerUserRolesDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user role id", http.StatusBadRequest)
		return
	}

	userRolesResponse, err := userRolesService.GetOne(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve user roles", err)
		return
	}

	lib.Success(w, "user roles successfuly retreived", userRolesResponse)
}

func HandlerUserRolesCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	var input api.UserRolesParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userRolesResponse, err := userRolesService.CreateUserRoles(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create user roles", err)
		return
	}

	lib.Success(w, "user roles successfully created", userRolesResponse)
}

func HandlerUserRolesUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user role id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	var input api.UserRolesParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userRolesResponse, err := userRolesService.UpdateUserRole(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update user roles", err)
		return
	}

	lib.Success(w, "user roles successfully updated", userRolesResponse)
}
