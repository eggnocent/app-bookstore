package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerAuthorsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authorsResponse, err := authorsService.List(ctx)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "failed to retrieve authors list", err)
		return
	}

	lib.Success(w, "success to retrieve authors list", authorsResponse)
}

func HandlerAuthorsDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid authors id", http.StatusBadRequest)
		return
	}

	authorsResponse, err := authorsService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "failed to retreive detail on authors", err)
		return
	}

	lib.Success(w, "success to retrieve detail on authors", authorsResponse)
}

func HandlerAuthorsCreate(w http.ResponseWriter, r *http.Request) {
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

	var input api.AuthorsParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	authorsResponse, err := authorsService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create authors", err)
		return
	}

	lib.Success(w, "authors successfully created", authorsResponse)
}

func HandlerAuthorsUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid authors id", http.StatusBadRequest)
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

	var input api.AuthorsParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	authorsResponse, err := authorsService.Update(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update authors", err)
		return
	}

	lib.Success(w, "authors successfully updated", authorsResponse)
}
