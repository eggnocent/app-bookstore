package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"app-bookstore/model"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerAuthorsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := lib.ParseQueryParam(ctx, r)
	if err != nil {
		http.Error(w, "invalid query parameter", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var dateFilter model.DateFilter
	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "invalid start_date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "invalid end_date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		dateFilter = model.DateFilter{
			StartDate: startDate,
			EndDate:   endDate,
		}
	}

	authorsResponse, err := authorsService.List(ctx, res, dateFilter)
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
