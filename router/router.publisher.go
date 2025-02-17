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

func HandlerPublisherList(w http.ResponseWriter, r *http.Request) {
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
	}

	publisherResponse, err := publisherService.List(ctx, res, dateFilter)
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
