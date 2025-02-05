package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"

	"github.com/jmoiron/sqlx"
)

var (
	userService *api.UserModule
	roleService *api.RoleModule
)

func Init(db *sqlx.DB, jwt lib.Jwt) {
	userService = api.NewUserModule(db, jwt)
	roleService = api.NewRoleModule(db, jwt)
}
