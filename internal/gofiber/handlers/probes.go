package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func LiveProbe(version string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"code":    http.StatusOK,
			"status":  http.StatusText(http.StatusOK),
			"version": version,
		})
	}
}
