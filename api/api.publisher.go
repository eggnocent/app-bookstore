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

type PublisherModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewPublisherModule(db *sqlx.DB, jwt lib.Jwt) *PublisherModule {
	return &PublisherModule{
		db:   db,
		name: "publisher-module",
		JWT:  jwt,
	}
}

type PublisherParam struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func (p *PublisherModule) List(ctx context.Context, filter lib.Filter, dateFilter model.DateFilter) ([]model.PublisherResponse, error) {
	publiserRequest, err := model.GetAllPublisher(ctx, p.db, filter, dateFilter)
	if err != nil {
		return nil, err
	}

	var response []model.PublisherResponse
	for _, publisher := range publiserRequest {
		response = append(response, publisher.Response())
	}

	return response, err
}

func (p *PublisherModule) Detail(ctx context.Context, id uuid.UUID) (model.PublisherResponse, error) {
	publisherRequest, err := model.GetOnePublisher(ctx, p.db, id)
	if err != nil {
		return model.PublisherResponse{}, err
	}

	return publisherRequest.Response(), nil
}

func (p *PublisherModule) Create(ctx context.Context, token string, param PublisherParam) (interface{}, error) {
	claims, err := p.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	publisher := model.PublisherModel{
		ID:        uuid.New(),
		Name:      param.Name,
		Address:   param.Address,
		Phone:     param.Phone,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = publisher.Insert(ctx, p.db)
	if err != nil {
		return nil, err
	}

	return publisher.Response(), nil
}

func (p *PublisherModule) Update(ctx context.Context, token string, param PublisherParam, id uuid.UUID) (interface{}, error) {
	claims, err := p.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	publisher := model.PublisherModel{
		ID:      id,
		Name:    param.Name,
		Address: param.Address,
		Phone:   param.Phone,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = publisher.Update(ctx, p.db)
	if err != nil {
		return nil, err
	}

	return publisher.Response(), nil
}
