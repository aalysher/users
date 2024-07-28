package validator

import (
	"fmt"
	"regexp"

	"users/internal/models"
)

func ValidateUser(user *models.User) error {
	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if !isValidEmail(user.Email) {
		return fmt.Errorf("invalid email address")
	}

	return nil
}

func ValidateUserUpdate(updates *models.UserUpdate) error {
	if updates.FirstName != nil && *updates.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if updates.LastName != nil && *updates.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if updates.Email != nil && !isValidEmail(*updates.Email) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
