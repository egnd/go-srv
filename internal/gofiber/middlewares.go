package gofiber

import (
	"strings"

	"github.com/egnd/go-toolbox/metrics/victoria"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/egnd/go-srv/internal/gofiber/middlewares"
)

func (srv *goFiber) setMiddlewares() *goFiber {
	srv.server.Use(
		favicon.New(),
		middlewares.RequestID(),
		middlewares.Logging(srv.logger),
		middlewares.HTTPMetrics(victoria.Opts{Namespace: "gosrv", Subsystem: "gofiber"},
			func(ctx *fiber.Ctx) bool {
				switch reqURI := string(ctx.Request().RequestURI()); {
				case strings.HasPrefix(reqURI, "/metrics"):
					return true
				case strings.HasPrefix(reqURI, "/debug/"):
					return true
				default:
					return false
				}
			},
		),
		recover.New(recover.Config{EnableStackTrace: srv.debug}),
	)

	srv.server.Use("/debug/pprof/*", pprof.New())

	return srv
}
