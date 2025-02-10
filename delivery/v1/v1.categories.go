package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPICategories(r *mux.Router) {
	r.HandleFunc("/categories", router.HandlerCategoriesList).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", router.HandlerCategoryDetail).Methods(http.MethodGet)
	r.HandleFunc("/categories", router.HandlerCategoryCreate).Methods(http.MethodPost)
	r.HandleFunc("/categories/{id}", router.HandlerCategoryUpdate).Methods(http.MethodPut)
}
