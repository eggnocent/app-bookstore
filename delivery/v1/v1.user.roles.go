package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIUserRoles(r *mux.Router) {
	r.HandleFunc("/user-roles", router.HandlerUserRole).Methods(http.MethodGet)
	r.HandleFunc("/user-roles/{id}", router.HandlerUserRolesDetail).Methods(http.MethodGet)
	r.HandleFunc("/user-roles", router.HandlerUserRolesCreate).Methods(http.MethodPost)
	r.HandleFunc("/user-roles/{id}", router.HandlerUserRolesUpdate).Methods(http.MethodPut)
}
