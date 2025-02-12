package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPILoans(r *mux.Router) {
	r.HandleFunc("/loans", router.HandlerListLoan).Methods(http.MethodGet)
	r.HandleFunc("/loans/{id}", router.HandlerLoansDetail).Methods(http.MethodGet)
	r.HandleFunc("/loans", router.HandlerLoanCreate).Methods(http.MethodPost)
	r.HandleFunc("/loans/{id}", router.HandlerLoansUpdate).Methods(http.MethodPut)
}
