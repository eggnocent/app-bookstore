package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIAuthors(r *mux.Router) {
	r.HandleFunc("/authors", router.HandlerAuthorsList).Methods(http.MethodGet)
	r.HandleFunc("/authors/{id}", router.HandlerAuthorsDetail).Methods(http.MethodGet)
	r.HandleFunc("/authors", router.HandlerAuthorsCreate).Methods(http.MethodPost)
	r.HandleFunc("/authors/{id}", router.HandlerAuthorsUpdate).Methods(http.MethodPut)
}
