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
		pprof.New(),
		middlewares.RequestID(),
		middlewares.Logging(srv.logger),
		middlewares.Metrics(victoria.Opts{Namespace: "gosrv", Subsystem: "gofiber"},
			func(ctx *fiber.Ctx) (skip bool) {
				reqURI := string(ctx.Request().RequestURI())
				switch {
				case strings.HasPrefix(reqURI, "/metrics"):
					skip = true
				case strings.HasPrefix(reqURI, "/debug/"):
					skip = true
				}

				return
			},
		),
		recover.New(recover.Config{EnableStackTrace: srv.debug}),
	)

	return srv
}
