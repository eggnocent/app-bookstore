package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIUser(r *mux.Router) {
	r.HandleFunc("/register", router.HandlerRegisterUser).Methods(http.MethodPost)
	r.HandleFunc("/login", router.HandlerLogin).Methods(http.MethodPost)
	r.HandleFunc("/change-password", router.HandlerChangePassword).Methods(http.MethodPost)
}
