package handlers

import (
	"github.com/egnd/go-srv/internal/ctxtools"
	"github.com/gofiber/fiber/v2"
)

func HelloWorld() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctxtools.GetLogger(ctx.UserContext()).Info("handle hello world request")

		return ctx.SendString("Hello, World!")
	}
}
