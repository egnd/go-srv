package middlewares

import (
	"time"

	"github.com/egnd/go-srv/internal/ctxtools"
	"github.com/go-logr/logr"
	"github.com/gofiber/fiber/v2"
)

func Logging(logger logr.Logger, skip ...SkipFunc) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if len(skip) > 0 && skip[0] != nil && skip[0](ctx) {
			return ctx.Next()
		}

		logger := logger.WithValues("req_id", ctxtools.GetReqID(ctx.UserContext()))

		ctx.SetUserContext(ctxtools.AddLogger(ctx.UserContext(), logger))

		logger.V(1).Info("request",
			"method", ctx.Method(),
			"uri", ctx.Request().RequestURI(),
			"size", len(ctx.Request().Body()),
		)

		start := time.Now()
		err := ctx.Next()

		logger.V(1).Info("response",
			"code", ctx.Response().StatusCode(),
			"dur_ms", time.Since(start).Milliseconds(),
			"size", len(ctx.Response().Body()),
		)

		return err
	}
}
