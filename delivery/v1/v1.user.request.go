package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIUserRequest(r *mux.Router) {
	r.HandleFunc("/user-request", router.HandlerUserRequestList).Methods(http.MethodGet)
	r.HandleFunc("/user-request/{id}", router.HandlerUserRequestDetail).Methods(http.MethodGet)
	r.HandleFunc("/user-request", router.HandlerUserRequest).Methods(http.MethodPost)
	r.HandleFunc("/user-request/{id}", router.HandlerUserRequestUpdate).Methods(http.MethodPut)
}
