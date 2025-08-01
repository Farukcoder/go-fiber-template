package middleware

import (
	"time"

	"garma_track/helpers"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger middleware logs all HTTP requests
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		helpers.LogRequest(
			c.Method(),
			c.Path(),
			c.IP(),
			c.Get("User-Agent"),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
