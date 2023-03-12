package middlewares

import (
	"time"

	"github.com/egnd/go-toolbox/metrics/victoria"
	"github.com/gofiber/fiber/v2"
)

func Metrics(metricsOpt victoria.Opts, skip ...SkipFunc) fiber.Handler {
	httpReqsOpt := metricsOpt
	httpReqsOpt.Name = "req_cnt"
	httpReqs := victoria.NewIncrement(&httpReqsOpt, "method", "uri")

	httpRespSizeOpt := metricsOpt
	httpRespSizeOpt.Name = "resp_bytes"
	httpRespSize := victoria.NewIncrement(&httpRespSizeOpt, "method", "uri")

	httpRespOpt := metricsOpt
	httpRespOpt.Name = "resp"
	httpResp := victoria.NewHisto(&httpRespOpt, "method", "uri")

	return func(ctx *fiber.Ctx) error {
		if len(skip) > 0 && skip[0] != nil && skip[0](ctx) {
			return ctx.Next()
		}

		httpReqs.With("method", ctx.Method(), "uri", string(ctx.Request().URI().Path())).Inc()

		start := time.Now()

		defer func() {
			httpResp.With(
				"method", ctx.Method(), "uri", string(ctx.Request().URI().Path()),
			).Update(time.Since(start).Seconds())
			httpRespSize.With(
				"method", ctx.Method(), "uri", string(ctx.Request().URI().Path()),
			).Add(len(ctx.Response().Body()))
		}()

		return ctx.Next()
	}
}