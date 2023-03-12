package handlers

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gofiber/fiber/v2"
)

func PromMetrics(ctx *fiber.Ctx) error {
	metrics.WriteProcessMetrics(ctx.Response().BodyWriter())
	metrics.WriteFDMetrics(ctx.Response().BodyWriter())
	metrics.WritePrometheus(ctx.Response().BodyWriter(), false)
	return nil
}
