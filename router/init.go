package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"

	"github.com/jmoiron/sqlx"
)

var (
	userService *api.UserModule
)

func Init(db *sqlx.DB, jwt lib.Jwt) {
	if db == nil {
		panic("❌ ERROR: router.Init received nil database connection")
	}
	if jwt == nil {
		panic("❌ ERROR: router.Init received nil JWT service")
	}

	userService = api.NewUserModule(db, jwt)
}
