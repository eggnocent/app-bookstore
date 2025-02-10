package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIPublisher(r *mux.Router) {
	r.HandleFunc("/publisher", router.HandlerPublisherList).Methods(http.MethodGet)
	r.HandleFunc("/publisher/{id}", router.HandlerPublishDetail).Methods(http.MethodGet)
	r.HandleFunc("/publisher", router.HandlerPublisherCreate).Methods(http.MethodPost)
	r.HandleFunc("/publisher/{id}", router.HandlerPubliserUpdate).Methods(http.MethodPut)
}
