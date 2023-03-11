package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

type SkipFunc func(c *fiber.Ctx) bool
