package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"
	"app-bookstore/model"
	"net/http"
	"time"
)

func HandlerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := lib.ParseQueryParam(ctx, r)
	if err != nil {
		http.Error(w, "invalid query parameter", http.StatusBadRequest)
		return
	}

	res.IsActive = r.URL.Query().Get("is_active") == "true"
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

	userResponse, err := userService.List(ctx, res, dateFilter)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "failed to retrieve user", err)
		return
	}

	lib.Success(w, "success to retrieve user", userResponse)
}

func HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input api.UserParam
	err := lib.ParseBody(ctx, r, &input)
	if err != nil {
		http.Error(w, "Failed to parse user"+err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := userService.Register(ctx, input)
	if err != nil {
		lib.Error(w, http.StatusBadRequest, "failed to register user", err)
		return
	}

	lib.Success(w, "user register successfully", userResponse)
}
