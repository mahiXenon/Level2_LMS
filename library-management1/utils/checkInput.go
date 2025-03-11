package utils

import (
	"regexp"
	"unicode"
)

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// isValidContactNumber checks if the contact number contains only digits
func IsValidContactNumber(contact string) bool {
	if len(contact) != 10 {
		return false
	}
	count := 0
	for i := 0; i < 10; i++ {
		if contact[i] >= '0' && contact[i] <= '9' {
			count = count + 1
		} else {
			break
		}
	}
	if count == 10 {
		return true
	}
	return false
}

// isValidPassword checks if password meets length and special character requirements
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasSpecial := false
	for _, ch := range password {
		if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
			hasSpecial = true
		}
	}
	return hasSpecial
}
