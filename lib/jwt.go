package lib

import (
	"app-bookstore/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtData struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Jwt interface {
	GenerateToken(data *JwtData) (string, int64, error)
	VerifyAccessToken(token string) (*JwtData, error)
}

type Options struct {
	SigningKey string
	Issuer     string
}

func (o *Options) GenerateToken(data *JwtData) (string, int64, error) {
	if data.UserID == "" {
		return "", 0, errors.New("user id is required")
	}

	now := time.Now().UTC()
	expiredAt := now.Add(24 * time.Hour)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiredAt)
	data.RegisteredClaims.Issuer = o.Issuer
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	signedToken, err := token.SignedString([]byte(o.SigningKey))
	if err != nil {
		return "", 0, err
	}
	return signedToken, expiredAt.Unix(), nil
}

func (o *Options) VerifyAccessToken(token string) (*JwtData, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JwtData{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(o.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*JwtData); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid or expired token")
}

func NewJWT(cfg *config.Config) Jwt {
	return &Options{
		SigningKey: cfg.App.JwtSecretKey,
		Issuer:     cfg.App.JwtIssuer,
	}
}
