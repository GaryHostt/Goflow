package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using go-playground/validator tags
// Returns nil if valid, or a user-friendly error message if invalid
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		// Convert validation errors to user-friendly messages
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// formatValidationErrors converts validator errors to user-friendly messages
func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string
	
	for _, err := range errs {
		field := err.Field()
		tag := err.Tag()
		
		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters", field, err.Param())
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
		case "url":
			message = fmt.Sprintf("%s must be a valid URL", field)
		case "oneof":
			message = fmt.Sprintf("%s must be one of: %s", field, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}
		
		messages = append(messages, message)
	}
	
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// Validate specific types with custom logic

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	type EmailValidator struct {
		Email string `validate:"required,email"`
	}
	
	return ValidateStruct(EmailValidator{Email: email})
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	
	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters")
	}
	
	return nil
}

// ValidateURL validates a URL
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL is required")
	}
	
	type URLValidator struct {
		URL string `validate:"required,url"`
	}
	
	return ValidateStruct(URLValidator{URL: url})
}

