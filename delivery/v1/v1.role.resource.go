package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIRoleResource(r *mux.Router) {
	r.HandleFunc("/role-resources", router.HandlerRoleResourceList).Methods(http.MethodGet)
	r.HandleFunc("/role-resources/{id}", router.HandlerRoleResourceDetail).Methods(http.MethodGet)
	r.HandleFunc("/role-resources", router.HandlerRoleResourceCreate).Methods(http.MethodPost)
	r.HandleFunc("/role-resources/{id}", router.HandlerRoleResourceUpdate).Methods(http.MethodPut)
}
