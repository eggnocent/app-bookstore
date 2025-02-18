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

type BookModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewBooksModule(db *sqlx.DB, jwt lib.Jwt) *BookModule {
	return &BookModule{
		db:   db,
		name: "books-module",
		JWT:  jwt,
	}
}

type BookParam struct {
	Title         string `json:"title"`
	AuthorID      string `json:"author_id"`
	PublisherID   string `json:"publisher_id"`
	CategoryID    string `json:"category_id"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
	Status        string `json:"status"`
	AccessLevel   string `json:"access_level"`
}

func (b *BookModule) List(ctx context.Context, filter lib.Filter, dateFilter model.DateFilter) ([]model.BookResponse, error) {
	bookRequest, err := model.GetAllBooks(ctx, b.db, filter, dateFilter)
	if err != nil {
		return nil, err
	}

	var response []model.BookResponse
	for _, book := range bookRequest {
		response = append(response, book.Response())
	}

	return response, nil
}

func (b *BookModule) Detail(ctx context.Context, id uuid.UUID) (model.BookResponse, error) {
	bookRequest, err := model.GetOneBooks(ctx, b.db, id)
	if err != nil {
		return model.BookResponse{}, err
	}
	return bookRequest.Response(), nil
}

func (b *BookModule) Create(ctx context.Context, token string, param BookParam) (interface{}, error) {
	claims, err := b.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	book := model.BookModel{
		ID:       uuid.New(),
		Title:    param.Title,
		AuthorID: uuid.MustParse(param.AuthorID),
		PublisherID: uuid.NullUUID{
			UUID:  uuid.MustParse(param.PublisherID),
			Valid: param.PublisherID != "",
		},
		CategoryID: uuid.NullUUID{
			UUID:  uuid.MustParse(param.CategoryID),
			Valid: param.CategoryID != "",
		},
		PublishedYear: param.PublishedYear,
		ISBN:          param.ISBN,
		Status:        param.Status,
		AccessLevel:   param.AccessLevel,
		CreatedAt:     time.Now(),
		CreatedBy:     userID,
	}

	err = book.Insert(ctx, b.db)
	if err != nil {
		return nil, err
	}

	return book.Response(), nil
}

func (b *BookModule) Update(ctx context.Context, token string, param BookParam, id uuid.UUID) (interface{}, error) {
	claims, err := b.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	book := model.BookModel{
		ID:       id,
		Title:    param.Title,
		AuthorID: uuid.MustParse(param.AuthorID),
		PublisherID: uuid.NullUUID{
			UUID:  uuid.MustParse(param.PublisherID),
			Valid: param.PublisherID != "",
		},
		CategoryID: uuid.NullUUID{
			UUID:  uuid.MustParse(param.CategoryID),
			Valid: param.CategoryID != "",
		},
		PublishedYear: param.PublishedYear,
		ISBN:          param.ISBN,
		Status:        param.Status,
		AccessLevel:   param.AccessLevel,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = book.Update(ctx, b.db, id)
	if err != nil {
		return nil, err
	}

	return book.Response(), nil
}
