package gofiber

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"

	"github.com/egnd/go-srv/internal/ctxtools"
	"github.com/egnd/go-srv/internal/gofiber/handlers"
	"github.com/egnd/go-srv/internal/gofiber/middlewares"
)

type goFiber struct {
	debug    bool
	disabled bool
	ctx      context.Context
	cfg      *viper.Viper
	logger   logr.Logger
	server   *fiber.App
}

func New(ctx context.Context, cfg *viper.Viper, logger logr.Logger) *goFiber {
	return (&goFiber{
		ctx:      ctx,
		cfg:      cfg,
		logger:   logger,
		debug:    cfg.GetBool("debug"),
		disabled: cfg.GetBool("disabled"),
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

func (srv *goFiber) initServer() *goFiber {
	srv.server = fiber.New(fiber.Config{
		Concurrency:           srv.cfg.GetInt("concurrency"),
		ReadBufferSize:        srv.cfg.GetInt("read_buffersize"),
		WriteBufferSize:       srv.cfg.GetInt("write_buffersize"),
		Prefork:               srv.cfg.GetBool("prefork"),
		ServerHeader:          srv.cfg.GetString("server_header"),
		StrictRouting:         srv.cfg.GetBool("routing.strict"),
		CaseSensitive:         srv.cfg.GetBool("routing.case_sensitive"),
		ReadTimeout:           srv.cfg.GetDuration("timeouts.read"),
		WriteTimeout:          srv.cfg.GetDuration("timeouts.write"),
		IdleTimeout:           srv.cfg.GetDuration("timeouts.idle"),
		GETOnly:               srv.cfg.GetBool("get_only"),
		DisableKeepalive:      srv.cfg.GetBool("disable_keepalive"),
		DisableStartupMessage: !srv.cfg.GetBool("startup_msg"),
		AppName:               srv.cfg.GetString("app_name"),
		StreamRequestBody:     srv.cfg.GetBool("stream_request_body"),
		ReduceMemoryUsage:     srv.cfg.GetBool("reduce_mem_usage"),
		EnablePrintRoutes:     srv.cfg.GetBool("print_routes"),
		ViewsLayout:           srv.cfg.GetString("views.layout"),
		RequestMethods:        srv.cfg.GetStringSlice("http_methods"),
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		ErrorHandler:          srv.errorHandler,
		// Views: ???,
	})

	return srv
}

func (srv *goFiber) errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	ctxtools.GetLogger(ctx.UserContext()).Error(err, "handler err", "http_code", code)
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return ctx.SendStatus(code)
}

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

func (srv *goFiber) setHandlers() *goFiber {
	srv.server.Get("/", handlers.HelloWorld())
	// @TODO: /live
	// @TODO: /metrics

	if srv.debug {
		srv.server.Get("/debug/dashboard", monitor.New(monitor.Config{Title: srv.server.Config().AppName}))
	}

	return srv
}
