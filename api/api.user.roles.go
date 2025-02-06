package api

import (
	"app-bookstore/lib"
	"app-bookstore/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRoleModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewUserRolesModule(db *sqlx.DB, jwt lib.Jwt) *UserRoleModule {
	return &UserRoleModule{
		db:   db,
		name: "user-roles-module",
		JWT:  jwt,
	}
}

type UserRolesParam struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}

func (uro *UserRoleModule) ListUserRole(ctx context.Context) ([]model.UserRoleResponse, error) {
	userRoles, err := model.GetAllUserRoles(ctx, uro.db)
	if err != nil {
		return nil, err
	}

	var response []model.UserRoleResponse
	for _, userRole := range userRoles {
		response = append(response, userRole.Response())
	}

	return response, nil
}

func (uro *UserRoleModule) GetOne(ctx context.Context, id uuid.UUID) (model.UserRoleResponse, error) {
	userRoles, err := model.GetOneUserRole(ctx, uro.db, id)
	if err != nil {
		return model.UserRoleResponse{}, err
	}

	return userRoles.Response(), nil
}

func (uro *UserRoleModule) CreateUserRoles(ctx context.Context, token string, param UserRolesParam) (interface{}, error) {
	claims, err := uro.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("failed to verify access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	userRolesID := uuid.New()
	userRoles := model.UserRoleModel{
		ID:        userRolesID,
		UserID:    param.UserID,
		RoleID:    param.RoleID,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = userRoles.Insert(ctx, uro.db)
	if err != nil {
		return nil, err
	}
	return userRoles.Response(), nil
}

func (uro *UserRoleModule) UpdateUserRole(ctx context.Context, token string, param UserRolesParam, id uuid.UUID) (interface{}, error) {
	claims, err := uro.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	userRole := model.UserRoleModel{
		ID:     id,
		UserID: param.UserID,
		RoleID: param.RoleID,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = userRole.Update(ctx, uro.db)
	if err != nil {
		return err, nil
	}

	return userRole.Response(), nil
}
