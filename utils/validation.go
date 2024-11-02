package utils

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}