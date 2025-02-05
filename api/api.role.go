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

type RoleModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewRoleModule(db *sqlx.DB, jwt lib.Jwt) *RoleModule {
	return &RoleModule{
		db:   db,
		name: "role-module",
		JWT:  jwt,
	}
}

type RoleParam struct {
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

func (r *RoleModule) List(ctx context.Context) ([]model.RoleResponse, error) {
	roles, err := model.GetAllRoles(ctx, r.db)
	if err != nil {
		return nil, err
	}

	var response []model.RoleResponse
	for _, role := range roles {
		response = append(response, role.Response())
	}

	return response, nil
}

func (r *RoleModule) Detail(ctx context.Context, id uuid.UUID) (model.RoleResponse, error) {
	roles, err := model.GetOneRole(ctx, r.db, id)
	if err != nil {
		return model.RoleResponse{}, err
	}

	return roles.Response(), nil
}

func (r *RoleModule) Create(ctx context.Context, token string, param RoleParam) (interface{}, error) {
	claims, err := r.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("failed to verify access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	roleID := uuid.New()
	role := model.RoleModel{
		ID:          roleID,
		Identifier:  param.Identifier,
		Description: param.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}

	err = role.Insert(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return role.Response(), nil
}

func (r *RoleModule) UpdateRole(ctx context.Context, token string, param RoleParam, id uuid.UUID) (interface{}, error) {
	claims, err := r.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	role := model.RoleModel{
		ID:          id,
		Identifier:  param.Identifier,
		Description: param.Description,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = role.Update(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return role.Response(), nil
}
