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

type RoleResourceModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewRoleResourceModule(db *sqlx.DB, jwt lib.Jwt) *RoleResourceModule {
	return &RoleResourceModule{
		db:   db,
		name: "role-resource-module",
		JWT:  jwt,
	}
}

type RoleResouceParam struct {
	RoleID     uuid.UUID `json:"role_id"`
	ResourceID uuid.UUID `json:"resource_id"`
	Method     string    `json:"method"`
	IsActive   bool      `json:"is_active"`
}

func (rr *RoleResourceModule) List(ctx context.Context) ([]model.RoleResourceResponse, error) {
	roleResources, err := model.GetAllRoleResource(ctx, rr.db)
	if err != nil {
		return nil, err
	}

	var response []model.RoleResourceResponse
	for _, roleResource := range roleResources {
		response = append(response, roleResource.Response())
	}

	return response, nil
}

func (rr *RoleResourceModule) Detail(ctx context.Context, id uuid.UUID) (model.RoleResourceResponse, error) {
	roleResource, err := model.GetOneRoleResource(ctx, rr.db, id)
	if err != nil {
		return model.RoleResourceResponse{}, err
	}

	return roleResource.Response(), nil
}

func (rr *RoleResourceModule) Create(ctx context.Context, token string, param RoleResouceParam) (interface{}, error) {
	claims, err := rr.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("failed to verify access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	roleResourceID := uuid.New()
	roleResource := model.RoleResourceModel{
		ID:         roleResourceID,
		RoleID:     param.RoleID,
		ResourceID: param.ResourceID,
		Method:     param.Method,
		IsActive:   param.IsActive,
		CreatedAt:  time.Now(),
		CreatedBy:  userID,
	}

	err = roleResource.Insert(ctx, rr.db)
	if err != nil {
		return nil, err
	}

	return roleResource.Response(), nil
}

func (rr *RoleResourceModule) Update(ctx context.Context, token string, param RoleResouceParam, id uuid.UUID) (interface{}, error) {
	claims, err := rr.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	roleResouce := model.RoleResourceModel{
		ID:         id,
		RoleID:     param.RoleID,
		ResourceID: param.ResourceID,
		Method:     param.Method,
		IsActive:   param.IsActive,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = roleResouce.Update(ctx, rr.db)
	if err != nil {
		return nil, err
	}

	return roleResouce.Response(), nil
}
