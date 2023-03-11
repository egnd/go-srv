package gofiber

import (
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
		recover.New(recover.Config{EnableStackTrace: srv.debug}),
	)

	return srv
}
