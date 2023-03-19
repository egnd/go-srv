package gofiber

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type goFiber struct {
	debug    bool
	disabled bool
	version  string
	ctx      context.Context
	cfg      *viper.Viper
	logger   logr.Logger
	server   *fiber.App
}

func New(ctx context.Context, cfg *viper.Viper, logger logr.Logger, version string) *goFiber {
	return (&goFiber{
		ctx:      ctx,
		cfg:      cfg,
		logger:   logger,
		debug:    cfg.GetBool("debug"),
		disabled: cfg.GetBool("disabled"),
		version:  version,
	}).initServer().setMiddlewares().setHandlers()
}

func (srv *goFiber) Start() error {
	if srv.disabled {
		srv.logger.Info("is disabled")

		return nil
	}

	srv.logger.Info("starting...", "port", srv.cfg.GetInt("port"))

	return srv.server.Listen(fmt.Sprintf(":%d", srv.cfg.GetInt("port")))
}

func (srv *goFiber) Stop() error {
	if srv.disabled {
		return nil
	}

	return srv.server.Shutdown()
}
