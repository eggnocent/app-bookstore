package lib

import (
	"context"
	"strings"
)

func SearchGenerate(ctx context.Context, operator string, search []string) string {
	if len(search) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(search, " "+operator+" ")
}
