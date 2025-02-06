package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIResource(r *mux.Router) {
	r.HandleFunc("/resource", router.HandlerResourceList).Methods(http.MethodGet)
	r.HandleFunc("/resource/{id}", router.HandlerResourceDetail).Methods(http.MethodGet)
	r.HandleFunc("/resource", router.HandlerResourceCreate).Methods(http.MethodPost)
	r.HandleFunc("/resource/{id}", router.HandlerResourceUpdate).Methods(http.MethodPut)
}
