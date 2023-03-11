package gofiber

import (
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/egnd/go-srv/internal/gofiber/handlers"
)

func (srv *goFiber) setHandlers() *goFiber {
	srv.server.Get("/", handlers.HelloWorld())
	// @TODO: /live
	// @TODO: /metrics

	if srv.debug {
		srv.server.Get("/debug/dashboard", monitor.New(monitor.Config{Title: srv.server.Config().AppName}))
	}

	return srv
}
