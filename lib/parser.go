package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Filter struct {
	Limit      int       `json:"limit" validate:"lte=100"`
	Offset     int       `json:"offset"`
	Dir        string    `json:"dir"`
	Search     string    `json:"search" validate:"omitempty,alphanum_space"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	CategoryID uuid.UUID `json:"category_id"`
	IsActive   bool      `json:"is_active"`
	IsPending  bool      `json:"is_pending"`
	IsApprove  bool      `json:"is_approve"`
	IsRejected bool      `json:"is_rejected"`
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("alphanum_space", isValidAlphanumWithSpace)
}

func isValidAlphanumWithSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	regex := regexp.MustCompile(`^[a-zA-Z0-9\s]*$`)
	return regex.MatchString(value)
}

func ParseQueryParam(ctx context.Context, r *http.Request) (Filter, error) {
	urisVal := r.URL.Query()
	var filter Filter

	limit, err := strconv.Atoi(urisVal.Get("limit"))
	if err == nil {
		filter.Limit = limit
	} else {
		filter.Limit = 10
	}

	offset, err := strconv.Atoi(urisVal.Get("offset"))
	if err == nil {
		filter.Offset = offset
	} else {
		filter.Offset = 0
	}

	dir := strings.ToUpper(urisVal.Get("dir"))
	if dir != "ASC" && dir != "DESC" {
		filter.Dir = "ASC"
	}
	filter.Dir = dir

	search := urisVal.Get("search")
	filter.Search = search

	if userIDStr := urisVal.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return Filter{}, errors.New("invalid user_id")
		}
		filter.UserID = userID
	}

	roleIDs := urisVal.Get("role_id")
	if roleIDs != "" {
		roleID, err := uuid.Parse(roleIDs)
		if err != nil {
			return Filter{}, errors.New("invalid request")
		}
		filter.RoleID = roleID
	}

	if categoryIDs := urisVal.Get("category_id"); categoryIDs != "" {
		categoryID, err := uuid.Parse(categoryIDs)
		if err != nil {
			return Filter{}, errors.New("invalid request")
		}
		filter.CategoryID = categoryID
	}

	filter.IsActive = urisVal.Get("is_active") == "true"
	filter.IsPending = urisVal.Get("is_pending") == "true"

	err = validate.Struct(filter)
	if err != nil {
		return Filter{}, errors.New("invalid request parameter")
	}

	return filter, nil
}

func ParseBody(ctx context.Context, r *http.Request, data interface{}) error {
	bData, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.Unmarshal(bData, &data)
	if err != nil {
		return err
	}

	err = validate.Struct(data)
	if err != nil {
		return err
	}

	return nil
}
