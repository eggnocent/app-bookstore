package lib

import (
	"strings"

	"github.com/google/uuid"
)

var SystemID = uuid.MustParse("c9c5f350-c0b2-4f09-bb9e-c1d134a371b1")

var (
	Active   = "ACTIVE"
	Pending  = "pending"
	Approve  = "approved"
	Reject   = "rejected"
	Inactive = "INACTIVE"
	Asc      = "ASC"
	Desc     = "DESC"

	StatusMap = map[string]string{
		"ACTIVE":   Active,
		"PENDING":  Pending,
		"APPROVED": Approve,
		"REJECTED": Reject,
		"INACTIVE": Inactive,
	}

	DirMap = map[string]string{
		"ASC":  Asc,
		"DESC": Desc,
	}
)

func IsValidStatus(status string) bool {
	_, exists := StatusMap[strings.ToUpper(status)]
	return exists
}

func GetValidStatus(status string) string {
	if valid, exists := StatusMap[strings.ToUpper(status)]; exists {
		return valid
	}
	return ""
}

func IsValidDirection(dir string) bool {
	_, exists := DirMap[strings.ToUpper(dir)]
	return exists
}

func GetValidDirection(dir string) string {
	if valid, exists := DirMap[strings.ToUpper(dir)]; exists {
		return valid
	}
	return "ASC" // Default ke ASC jika tidak valid
}
