package v1

import (
	"app-bookstore/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIResetPass(r *mux.Router) {
	r.HandleFunc("/request-reset-password", router.HandlerRequestResetPassword).Methods(http.MethodPost)
	r.HandleFunc("/validate-reset-token", router.HandlerValidateResetToken).Methods(http.MethodGet)
	r.HandleFunc("/reset-password", router.HandlerResetPassword).Methods(http.MethodPost)
}
