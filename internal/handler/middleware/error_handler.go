package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
