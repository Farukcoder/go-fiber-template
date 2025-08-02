package requests

import (
	"github.com/go-playground/validator/v10"
)

// LoginRequest represents the login request structure with validation rules
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents the registration request structure with validation rules
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	UserType string `json:"user_type,omitempty"`
}

// ValidateLoginRequest validates the login request
func (r *LoginRequest) Validate() []string {
	validate := validator.New()
	var errors []string

	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errors = append(errors, "Email is required and must be a valid email address")
			case "Password":
				errors = append(errors, "Password must be at least 6 characters long")
			}
		}
	}
	return errors
}

// ValidateRegisterRequest validates the registration request
func (r *RegisterRequest) Validate() []string {
	validate := validator.New()
	var errors []string

	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				errors = append(errors, "Name must be at least 2 characters long")
			case "Email":
				errors = append(errors, "Email is required and must be a valid email address")
			case "Password":
				errors = append(errors, "Password must be at least 6 characters long")
			}
		}
	}
	return errors
}
