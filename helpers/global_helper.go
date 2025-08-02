package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Response types for consistent API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// GenerateJWTToken generates a new JWT token for a user
func GenerateJWTToken(userID uint, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// SuccessResponse sends a success response with data
func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, status int, message string, errors []string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *fiber.Ctx, errors []string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", errors)
}

// LogRequest logs HTTP request details
func LogRequest(method, path, ip, userAgent string, status int, duration time.Duration) {
	// Create a more readable console output
	var statusColor string
	var statusIcon string

	switch {
	case status >= 200 && status < 300:
		statusColor = "✅"
		statusIcon = "SUCCESS"
	case status >= 300 && status < 400:
		statusColor = "🔄"
		statusIcon = "REDIRECT"
	case status >= 400 && status < 500:
		statusColor = "⚠️"
		statusIcon = "CLIENT_ERROR"
	case status >= 500:
		statusColor = "❌"
		statusIcon = "SERVER_ERROR"
	default:
		statusColor = "❓"
		statusIcon = "UNKNOWN"
	}

	// Console output with colors and formatting
	consoleMessage := fmt.Sprintf("%s [%s] %s %s - %d (%s) - %s - %v",
		statusColor,
		time.Now().Format("15:04:05"),
		method,
		path,
		status,
		statusIcon,
		ip,
		duration)

	// File output (original format)
	fileMessage := fmt.Sprintf("[HTTP] %s %s - %d - %s - %s - %v", method, path, status, ip, userAgent, duration)

	// Log to console with custom formatting
	fmt.Println(consoleMessage)

	// Also log to file using the existing logger
	if status >= 400 {
		Error(fileMessage, nil)
	} else {
		Info(fileMessage)
	}
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *fiber.Ctx) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", []string{"Invalid credentials"})
}

// ServerErrorResponse sends a server error response
func ServerErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, "Server error", []string{message})
}

// GetAuthUser extracts user claims from JWT token
func GetAuthUser(c *fiber.Ctx) jwt.MapClaims {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return nil
	}
	return claims
}

// ExtractBearerToken extracts the token from Authorization header
func ExtractBearerToken(c *fiber.Ctx) string {
	auth := c.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}
