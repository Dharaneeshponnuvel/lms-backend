package auth

import "regexp"

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[@$!%*?&]`)

	return lower.MatchString(password) &&
		upper.MatchString(password) &&
		number.MatchString(password) &&
		special.MatchString(password)
}

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
