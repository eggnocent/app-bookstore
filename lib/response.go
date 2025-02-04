package lib

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  *string     `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, status string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ApiResponse{
		Status: status,
		Data:   data,
	})
}

func Error(w http.ResponseWriter, statusCode int, status string, err error) {
	errMsg := err.Error()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ApiResponse{
		Status: status,
		Error:  &errMsg,
	})
}
