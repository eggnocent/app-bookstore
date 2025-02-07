package api

import (
	"app-bookstore/lib"
	"app-bookstore/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewUserModule(db *sqlx.DB, jwt lib.Jwt) *UserModule {
	return &UserModule{
		db:   db,
		name: "user-module",
		JWT:  jwt,
	}
}

type UserParam struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
}

type ChangePasswordParam struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (u *UserModule) List(ctx context.Context) ([]model.UserResponse, error) {
	users, err := model.GetAllUser(ctx, u.db)
	if err != nil {
		return nil, err
	}

	var response []model.UserResponse
	for _, user := range users {
		response = append(response, user.Response())
	}
	return response, nil
}

func (u *UserModule) Register(ctx context.Context, param UserParam) (interface{}, error) {
	hashPassword, err := lib.HashPassword(param.Password)
	if err != nil {
		return nil, err
	}

	user := &model.UserModel{
		ID:        uuid.New(),
		Username:  param.Username,
		Password:  hashPassword,
		CreatedAt: time.Now(),
		CreatedBy: lib.SystemID,
	}

	err = user.Insert(ctx, u.db)
	if err != nil {
		return nil, err
	}

	return user.Response(), nil
}

func (u *UserModule) Login(ctx context.Context, param UserParam) (*LoginResponse, error) {
	user, err := model.GetUserByUsername(ctx, u.db, param.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !lib.CheckPassword(param.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	tokenData := &lib.JwtData{
		UserID: user.ID.String(),
	}

	token, expiredAt, err := u.JWT.GenerateToken(tokenData)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	session := model.SessionModel{
		UserID:     user.ID,
		Token:      token,
		Expiration: time.Unix(expiredAt, 0),
		CreatedAt:  time.Now(),
		CreatedBy:  user.ID,
	}

	err = session.Insert(ctx, u.db)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		ExpiresAt:   time.Unix(expiredAt, 0),
	}, nil
}

func (u *UserModule) ChangePassword(ctx context.Context, token string, param ChangePasswordParam) error {
	claims, err := u.JWT.VerifyAccessToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return errors.New("invalid user id in token")
	}

	user, err := model.GetUserByID(ctx, u.db, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !lib.CheckPassword(param.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	if param.NewPassword != param.ConfirmPassword {
		return errors.New("new password and confirm password do not macth")
	}

	hashedPassword, err := lib.HashPassword(param.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = model.Update(ctx, u.db, userID, hashedPassword)
	if err != nil {
		return errors.New("failed to change password")
	}

	return nil
}

func (u *UserModule) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("token is invalid")
	}

	err := model.DeleteSessionByToken(ctx, u.db, token)
	if err != nil {
		return err
	}
	return nil
}
