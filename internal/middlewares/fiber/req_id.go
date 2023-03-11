package fiber

import (
	"github.com/egnd/go-srv/internal/ctxtools"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type RequestIDCfg struct {
	Header    string
	Generator func() string
	Next      SkipFunc
}

func RequestID(config ...RequestIDCfg) fiber.Handler {
	cfg := RequestIDCfg{
		Header:    "X-Request-ID",
		Generator: func() string { return ulid.Make().String() },
	}

	if len(config) > 0 && config[0].Header != "" {
		cfg.Header = config[0].Header
	}

	if len(config) > 0 && config[0].Generator != nil {
		cfg.Generator = config[0].Generator
	}

	return func(ctx *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(ctx) {
			return ctx.Next()
		}

		var rid string
		if rid = ctx.Get(cfg.Header); rid == "" {
			rid = cfg.Generator()
		}

		ctx.Set(cfg.Header, rid)
		ctx.SetUserContext(ctxtools.AddReqID(ctx.UserContext(), rid))

		return ctx.Next()
	}
}
