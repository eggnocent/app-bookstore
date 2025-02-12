package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerRatingList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ratingResponse, err := ratingService.List(ctx)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve rating list", err)
		return
	}

	lib.Success(w, "success to retrieve rating list", ratingResponse)
}

func HandlerRatingDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid rating id", http.StatusBadRequest)
		return
	}

	ratingResponse, err := ratingService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve rating detail", err)
		return
	}

	lib.Success(w, "success to retrieve rating detail", ratingResponse)
}

func HandlerRatingCreate(w http.ResponseWriter, r *http.Request) {
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

	var input api.RatingParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ratingResponse, err := ratingService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create rating", err)
		return
	}

	lib.Success(w, "rating successfully created", ratingResponse)
}

func HandlerRatingUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid rating id", http.StatusBadRequest)
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

	var input api.RatingUpdateParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body"+err.Error(), http.StatusInternalServerError)
		return
	}

	ratingResponse, err := ratingService.Update(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to update rating", err)
		return
	}

	lib.Success(w, "rating successfully updated", ratingResponse)
}
