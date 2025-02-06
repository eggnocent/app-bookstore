package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerRoleResourceList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roleResourceResponse, err := roleResourceService.List(ctx)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retreive role resource", err)
		return
	}

	lib.Success(w, "success to retreive role resource", roleResourceResponse)
}

func HandlerRoleResourceDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid resource id", http.StatusBadRequest)
		return
	}

	roleResourceResponse, err := roleResourceService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retreive role resource detail", err)
		return
	}

	lib.Success(w, "success to retreive detail on role resource", roleResourceResponse)
}

func HandlerRoleResourceCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "invalid token", http.StatusInternalServerError)
		return
	}

	var input api.RoleResouceParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	roleResourceResponse, err := roleResourceService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create role resource", err)
		return
	}

	lib.Success(w, "role resource created", roleResourceResponse)
}

func HandlerRoleResourceUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid resource id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "invalid token", http.StatusInternalServerError)
		return
	}

	var input api.RoleResouceParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	roleResourceResponse, err := roleResourceService.Update(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update role resource", err)
		return
	}

	lib.Success(w, "role resource successfully updated", roleResourceResponse)
}
