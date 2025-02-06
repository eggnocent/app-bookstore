package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerResourceList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resourceResponse, err := resourceService.List(ctx)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve resource", err)
		return
	}

	lib.Success(w, "success to retreive resource", resourceResponse)
}

func HandlerResourceDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid resource id", http.StatusBadRequest)
		return
	}

	resourceResponse, err := resourceService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve resource detail", err)
		return
	}

	lib.Success(w, "success to retreive detail on resource", resourceResponse)
}

func HandlerResourceCreate(w http.ResponseWriter, r *http.Request) {
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

	var input api.ResourceParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resourceResponse, err := resourceService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create resource", err)
		return
	}

	lib.Success(w, "resource successfully created", resourceResponse)
}

func HandlerResourceUpdate(w http.ResponseWriter, r *http.Request) {
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

	var input api.ResourceParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resourceResponse, err := resourceService.Update(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update resource", err)
		return
	}

	lib.Success(w, "resource successfully updated", resourceResponse)
}
