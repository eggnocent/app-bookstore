package lib

import "regexp"

func NormalizeEndpoint(url string) string {
	uuidPattern := regexp.MustCompile(`[a-f0-9-]{36}`)
	return uuidPattern.ReplaceAllString(url, "{id}")
}
