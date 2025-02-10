package api

import (
	"app-bookstore/lib"
	"app-bookstore/model"
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type PasswordResetModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewPasswordResetModule(db *sqlx.DB, jwt lib.Jwt) *PasswordResetModule {
	return &PasswordResetModule{
		db:   db,
		name: "password_reset",
		JWT:  jwt,
	}
}

type PasswordResetRequestParam struct {
	Username string `json:"username"`
}

type PasswordResetConfigParam struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (pr *PasswordResetModule) ValidateToken(ctx context.Context, token string) (interface{}, error) {
	userID, err := model.ValidateResetToken(ctx, pr.db, token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return map[string]string{
		"message": "Token is valid",
		"user_id": userID.String(),
	}, nil
}

func (pr *PasswordResetModule) RequestReset(ctx context.Context, param PasswordResetRequestParam) (interface{}, error) {
	user, err := model.GetUserByUsername(ctx, pr.db, param.Username)
	if err != nil {
		return nil, errors.New("username not found")
	}

	token, err := lib.GenerateResetToken()
	if err != nil {
		return nil, errors.New("failed to generate reset token")
	}

	passwordReset := model.PasswordResetModel{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}

	err = passwordReset.Insert(ctx, pr.db)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"message": "this token for testing",
		"token":   token,
	}, nil
}

func (pr *PasswordResetModule) ResetPassword(ctx context.Context, param PasswordResetConfigParam) (interface{}, error) {
	userID, err := model.ValidateResetToken(ctx, pr.db, param.Token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	hashedPass, err := lib.HashPassword(param.NewPassword)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	err = model.UpdatePassword(ctx, pr.db, userID, hashedPass)
	if err != nil {
		return nil, errors.New("failed to update password")
	}

	err = model.DeleteResetToken(ctx, pr.db, param.Token)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"message": "Password successfully reset",
	}, nil
}
