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

func HandlerListLoan(w http.ResponseWriter, r *http.Request) {
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

	loanResponse, err := loanService.List(ctx, res, dateFilter)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve loan list", err)
		return
	}

	lib.Success(w, "success to retrieve loan list", loanResponse)
}

func HandlerLoansDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid loan id", http.StatusBadRequest)
		return
	}

	loanResponse, err := loanService.Detail(ctx, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to retrieve loan detail", err)
		return
	}

	lib.Success(w, "success to retrieve loan detail", loanResponse)
}

func HandlerLoanCreate(w http.ResponseWriter, r *http.Request) {
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

	var input api.LoansParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	loanResponse, err := loanService.Create(ctx, token, input)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to create loan", err)
		return
	}

	lib.Success(w, "loan successfully created", loanResponse)
}

func HandlerLoansUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid loan id", http.StatusBadRequest)
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

	var input api.ReturnParam
	err = lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	loanResponse, err := loanService.Return(ctx, token, input, id)
	if err != nil {
		lib.Error(w, http.StatusInternalServerError, "failed to return loan", err)
		return
	}

	lib.Success(w, "loan successfully returned", loanResponse)
}
