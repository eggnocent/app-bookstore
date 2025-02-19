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

type LoansModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewLoansModule(db *sqlx.DB, jwt lib.Jwt) *LoansModule {
	return &LoansModule{
		db:   db,
		name: "loans-module",
		JWT:  jwt,
	}
}

type LoansParam struct {
	BookID     uuid.UUID `json:"book_id"`
	MemberID   uuid.UUID `json:"member_id"`
	LoanDate   time.Time `json:"loan_date"`
	ReturnDate time.Time `json:"return_date"`
	Status     string    `json:"status"`
}

type ReturnParam struct {
	Status string `json:"status"`
}

func (l *LoansModule) List(ctx context.Context, filter lib.Filter, dateFilter model.DateFilter) ([]model.LoansResponse, error) {
	loanResponse, err := model.GetAllLoans(ctx, l.db, filter, dateFilter)
	if err != nil {
		return nil, err
	}

	var response []model.LoansResponse
	for _, loan := range loanResponse {
		response = append(response, loan.Response())
	}

	return response, nil
}

func (l *LoansModule) Detail(ctx context.Context, id uuid.UUID) (model.LoansResponse, error) {
	loanResponse, err := model.GetOneLoans(ctx, l.db, id)
	if err != nil {
		return model.LoansResponse{}, err
	}

	return loanResponse.Response(), nil
}

func (l *LoansModule) Create(ctx context.Context, token string, param LoansParam) (interface{}, error) {
	claims, err := l.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	loans := model.LoansModel{
		ID:       uuid.New(),
		BookID:   param.BookID,
		MemberID: param.MemberID,
		LoanDate: time.Now(),
		ReturnDate: pq.NullTime{
			Time:  param.ReturnDate,
			Valid: true,
		},
		Status:    "borrowed",
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = loans.Insert(ctx, l.db)
	if err != nil {
		return nil, err
	}

	return loans.Response(), nil
}

func (l *LoansModule) Return(ctx context.Context, token string, param ReturnParam, id uuid.UUID) (interface{}, error) {
	claims, err := l.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	var returnDate pq.NullTime
	if param.Status == "returned" {
		returnDate = pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	loans := model.LoansModel{
		ID:         id,
		Status:     param.Status,
		ReturnDate: returnDate,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = loans.Update(ctx, l.db, id)
	if err != nil {
		return nil, err
	}
	return loans.Response(), nil
}
