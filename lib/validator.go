package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid validation error: " + err.Error())
		}

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errors.New("unexpected error during validation")
		}

		var errorMessage []string
		for _, validationErr := range validationErrors {
			switch validationErr.Tag() {
			case "required":
				errorMessage = append(errorMessage, fmt.Sprintf("Field %s is required", validationErr.Field()))
			case "email":
				errorMessage = append(errorMessage, fmt.Sprintf("Field %s must be a valid email", validationErr.Field()))
			default:
				errorMessage = append(errorMessage, fmt.Sprintf("Field %s is invalid", validationErr.Field()))
			}
		}

		return errors.New("Validation failed: " + joinMessage(errorMessage))
	}
	return nil
}

func ParseBody(ctx context.Context, r *http.Request, data interface{}) error {
	// Membaca body request
	bData, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body request: %w", err)
	}
	defer r.Body.Close()

	// Memastikan body tidak kosong
	if len(bData) == 0 {
		return errors.New("request body cannot be empty")
	}

	// Unmarshal JSON ke pointer struct
	err = json.Unmarshal(bData, &data) // Memastikan data berupa pointer
	if err != nil {
		return fmt.Errorf("failed to process JSON: %w", err)
	}

	// Validasi data
	err = ValidateStruct(data)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func joinMessage(messages []string) string {
	result := ""
	for i, message := range messages {
		if i > 0 {
			result += ", "
		}
		result += message
	}
	return result
}
