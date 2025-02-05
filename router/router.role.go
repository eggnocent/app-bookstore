package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerRoleGetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roleResponse, err := roleService.List(ctx)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "failed to retreive role", err)
		return
	}

	lib.Success(w, "success to retreive role", roleResponse)
}

func HandlerRolesDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid role id", http.StatusBadRequest)
		return
	}

	roleResponse, err := roleService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "faoled tp retrieve detail on roles", err)
		return
	}

	lib.Success(w, "success to retrieve detail on roles", roleResponse)
}

func HandlerRoleCreate(w http.ResponseWriter, r *http.Request) {
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

	var input api.RoleParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	roleResponse, err := roleService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create role", err)
		return
	}

	lib.Success(w, "role successfully created", roleResponse)

}

func HandlerRoleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid role id", http.StatusBadRequest)
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

	var input api.RoleParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	roleResponse, err := roleService.UpdateRole(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update role", err)
		return
	}

	lib.Success(w, "role successuly updated", roleResponse)
}
