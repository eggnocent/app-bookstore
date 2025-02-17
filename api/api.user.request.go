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

type UserRequestModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewUserRequestModule(db *sqlx.DB, jwt lib.Jwt) *UserRequestModule {
	return &UserRequestModule{
		db:   db,
		name: "user-request-module",
		JWT:  jwt,
	}
}

type UserRequestParam struct {
	UserID        uuid.UUID `json:"user_id"`
	RequestRoleID uuid.UUID `json:"requested_role_id"`
}

type UserRequestUpdateParam struct {
	Status string `json:"status"`
}

func (ur *UserRequestModule) ListUserRequest(ctx context.Context) ([]model.UserRequestResponse, error) {
	userRequest, err := model.GetAllUserRequest(ctx, ur.db)
	if err != nil {
		return nil, err
	}

	var response []model.UserRequestResponse
	for _, userRequest := range userRequest {
		response = append(response, userRequest.Response())
	}
	return response, err
}

func (ur *UserRequestModule) GetOneUserRequest(ctx context.Context, id uuid.UUID) (model.UserRequestResponse, error) {
	userRequest, err := model.GetOneUserRequest(ctx, ur.db, id)
	if err != nil {
		return model.UserRequestResponse{}, err
	}

	return userRequest.Response(), nil

}

func (ur *UserRequestModule) CreateUserRequest(ctx context.Context, token string, param UserRequestParam) (interface{}, error) {
	claims, err := ur.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("failed to verify access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	roleExists, err := model.CheckRoleExist(ctx, ur.db, param.RequestRoleID)
	if err != nil {
		return nil, errors.New("failed to check role existence")
	}

	if !roleExists {
		return nil, errors.New("requested role does not exist")
	}

	userRequest := model.UserRequestModel{
		ID:              uuid.New(),
		UserID:          param.UserID,
		RequestUserRole: param.RequestRoleID,
		Status:          "pending",
		CreatedAt:       time.Now(),
		CreatedBy:       userID,
	}

	err = userRequest.Insert(ctx, ur.db)
	if err != nil {
		return nil, err
	}

	return userRequest.Response(), nil
}

func (ur *UserRequestModule) UpdateUserRequest(ctx context.Context, token string, param UserRequestUpdateParam, id uuid.UUID) (interface{}, error) {
	claims, err := ur.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	userRequest := model.UserRequestModel{
		ID:     id,
		Status: param.Status,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = userRequest.Update(ctx, ur.db, id)
	if err != nil {
		return nil, err
	}

	return userRequest.Response(), nil
}
