package router

import (
	"app-bookstore/api"
	"app-bookstore/lib"

	"github.com/jmoiron/sqlx"
)

var (
	userService         *api.UserModule
	roleService         *api.RoleModule
	userRequestService  *api.UserRequestModule
	userRolesService    *api.UserRoleModule
	resourceService     *api.ResourceModule
	roleResourceService *api.RoleResourceModule
)

func Init(db *sqlx.DB, jwt lib.Jwt) {
	userService = api.NewUserModule(db, jwt)
	roleService = api.NewRoleModule(db, jwt)
	userRequestService = api.NewUserRequestModule(db, jwt)
	userRolesService = api.NewUserRolesModule(db, jwt)
	resourceService = api.NewResourceModule(db, jwt)
	roleResourceService = api.NewRoleResourceModule(db, jwt)
}
