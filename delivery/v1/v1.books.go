package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIBooks(r *mux.Router) {
	r.HandleFunc("/books", router.HandlerBookList).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", router.HandlerBookDetail).Methods(http.MethodGet)
	r.HandleFunc("/books", router.HandlerBookCreate).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", router.HandlerBookUpdate).Methods(http.MethodPut)
}
