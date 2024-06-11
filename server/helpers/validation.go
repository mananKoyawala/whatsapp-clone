package helper

import (
	"regexp"
	"strconv"
)

func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func IsNonNegative(n int) bool {
	return n > 0
}

func ValidatePassword(password string) (string, bool) {
	minLength := 6
	maxLength := 10
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecial := false

	if password == "" {
		return "Password is required", false
	}

	if len(password) < minLength || len(password) > maxLength {
		return "Password must be between 6 to 10 characters", false
	}

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLowercase = true
		case 'A' <= char && char <= 'Z':
			hasUppercase = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case char == '@' || char == '$' || char == '!' || char == '*' || char == '%' || char == '?' || char == '&':
			hasSpecial = true
		}
	}

	if !hasLowercase {
		return "Password must contain at least one lowercase letter", false
	}

	if !hasUppercase {
		return "Password must contain at least one uppercase letter", false
	}

	if !hasDigit {
		return "Password must contain at least one digit", false
	}

	if !hasSpecial {
		return "Password must contain at least one special character (@$!*%?&)", false
	}

	return "", true // Password passed all validations
}

func CheckLength(num int, length int) bool {
	return len(strconv.Itoa(num)) != length
}
