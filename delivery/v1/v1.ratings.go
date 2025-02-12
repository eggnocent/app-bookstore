package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIRating(r *mux.Router) {
	r.HandleFunc("/ratings", router.HandlerRatingList).Methods(http.MethodGet)
	r.HandleFunc("/ratings/{id}", router.HandlerRatingDetail).Methods(http.MethodGet)
	r.HandleFunc("/ratings", router.HandlerRatingCreate).Methods(http.MethodPost)
	r.HandleFunc("/ratings/{id}", router.HandlerRatingUpdate).Methods(http.MethodPut)
}
