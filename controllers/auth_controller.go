package controllers

import (
	"go-fiber-template/helpers"
	"go-fiber-template/models"
	"go-fiber-template/requests"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	input := new(requests.LoginRequest)
	if err := c.BodyParser(input); err != nil {
		helpers.Error("Failed to parse login input: %v", err)
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	if errors := input.Validate(); len(errors) > 0 {
		return helpers.ValidationErrorResponse(c, errors)
	}

	var user models.User
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return helpers.UnauthorizedResponse(c)
	}

	if err := user.ComparePassword(input.Password); err != nil {
		return helpers.UnauthorizedResponse(c)
	}

	tokenString, err := helpers.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		return helpers.ServerErrorResponse(c, "Could not generate token")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
		"token": tokenString,
		"user":  user,
	})
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	input := new(requests.RegisterRequest)
	if err := c.BodyParser(input); err != nil {
		helpers.Error("Failed to parse registration input: %v", err)
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	if errors := input.Validate(); len(errors) > 0 {
		return helpers.ValidationErrorResponse(c, errors)
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		UserType: models.Employee, // Default user type
	}

	// Set user type if provided
	if input.UserType != "" {
		switch input.UserType {
		case "system_admin":
			user.UserType = models.SystemAdmin
		case "garments_admin":
			user.UserType = models.GarmentsAdmin
		case "department_admin":
			user.UserType = models.DepartmentAdmin
		case "employee":
			user.UserType = models.Employee
		default:
			user.UserType = models.Employee
		}
	}

	if err := user.HashPassword(); err != nil {
		return helpers.ServerErrorResponse(c, "Could not hash password")
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		return helpers.ServerErrorResponse(c, "Could not create user")
	}

	user.Password = "" // Don't send password in response
	return helpers.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", user)
}
