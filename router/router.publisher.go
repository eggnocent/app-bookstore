package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerPublisherList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	publisherResponse, err := publisherService.List(ctx)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve publisher", err)
		return
	}

	lib.Success(w, "success to retreive publisher list", publisherResponse)
}

func HandlerPublishDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid publisher id", http.StatusBadRequest)
		return
	}

	publisherResponse, err := publisherService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve publisher detail", err)
		return
	}

	lib.Success(w, "success to retrieve publisher detail", publisherResponse)
}

func HandlerPublisherCreate(w http.ResponseWriter, r *http.Request) {
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

	var input api.PublisherParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	publisherResponse, err := publisherService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create publisher", err)
		return
	}

	lib.Success(w, "publisher successfully created", publisherResponse)
}

func HandlerPubliserUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid publisher id", http.StatusBadRequest)
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

	var input api.PublisherParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	publisherResponse, err := publisherService.Update(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update publisher", err)
		return
	}

	lib.Success(w, "publisher successfully updated", publisherResponse)
}
