package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIRole(r *mux.Router) {
	r.HandleFunc("/roles", router.HandlerRoleGetAll).Methods(http.MethodGet)
	r.HandleFunc("/roles/{id}", router.HandlerRolesDetail).Methods(http.MethodGet)
	r.HandleFunc("/roles", router.HandlerRoleCreate).Methods(http.MethodPost)
	r.HandleFunc("/roles/{id}", router.HandlerRoleUpdate).Methods(http.MethodPut)

}
