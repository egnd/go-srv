package gofiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/egnd/go-srv/internal/ctxtools"
)

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
