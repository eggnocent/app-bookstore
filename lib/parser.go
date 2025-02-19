package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Filter struct {
	Limit         int       `json:"limit" validate:"lte=100"`
	Offset        int       `json:"offset"`
	Dir           string    `json:"dir"`
	Search        string    `json:"search" validate:"omitempty,alphanum_space"`
	UserID        uuid.UUID `json:"user_id"`
	RoleID        uuid.UUID `json:"role_id"`
	CategoryID    uuid.UUID `json:"category_id"`
	IsPending     bool      `json:"is_pending"`
	IsApprove     bool      `json:"is_approve"`
	IsRejected    bool      `json:"is_rejected"`
	AuthorBook    string    `json:"author_book"`
	PublisherID   uuid.UUID `json:"publisher_id"`
	AuthorID      uuid.UUID `json:"author_id"`
	PublishedYear int       `json:"published_year"`
	Status        string    `json:"status"`
	AccessLevel   string    `json:"access_level"`
	Available     bool      `json:"available"`
	Borrowed      bool      `json:"borrowed"`
	Public        bool      `json:"public"`
	MemberOnly    bool      `json:"member_only"`
	AdminOnly     bool      `json:"admin_only"`
	MemberID      uuid.UUID `json:"member_id"`
	Returned      bool      `json:"returned"`
	BookID        uuid.UUID `json:"book_id"`
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

	if publisherIDs := urisVal.Get("publisher_id"); publisherIDs != "" {
		publisherID, err := uuid.Parse(publisherIDs)
		if err != nil {
			return Filter{}, errors.New("invalid request")
		}
		filter.PublisherID = publisherID
	}

	if authorIDs := urisVal.Get("author_id"); authorIDs != "" {
		authorID, err := uuid.Parse(authorIDs)
		if err != nil {
			return Filter{}, errors.New("invalid request")
		}
		filter.AuthorID = authorID
	}

	if memberID := urisVal.Get("member_id"); memberID != "" {
		memberID, err := uuid.Parse(memberID)
		if err != nil {
			return Filter{}, errors.New("invalid request")
		}
		filter.MemberID = memberID
	}

	log.Println(filter.MemberID)

	if publishedYear := urisVal.Get("published_year"); publishedYear != "" {
		publishedYearInt, err := strconv.Atoi(publishedYear)
		if err != nil {
			return Filter{}, errors.New("invalid published_year")
		}
		filter.PublishedYear = publishedYearInt
	}

	if booksID := urisVal.Get("book_id"); booksID != "" {
		bookID, err := uuid.Parse(booksID)
		if err != nil {
			return Filter{}, errors.New("invalid request")
		}
		filter.BookID = bookID
	}

	filter.IsPending = urisVal.Get("is_pending") == "true"
	filter.IsApprove = urisVal.Get("is_approve") == "true"
	filter.IsRejected = urisVal.Get("is_rejected") == "true"
	filter.Available = urisVal.Get("is_available") == "true"
	filter.Borrowed = urisVal.Get("is_borrowed") == "true"
	filter.Returned = urisVal.Get("is_returned") == "true"
	filter.Public = urisVal.Get("is_public") == "true"
	filter.MemberOnly = urisVal.Get("is_member_only") == "true"
	filter.AdminOnly = urisVal.Get("is_admin_only") == "true"

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
