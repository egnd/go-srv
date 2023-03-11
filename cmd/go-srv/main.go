package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/egnd/go-srv/internal/services"
	"github.com/egnd/go-toolbox/config"
	"github.com/egnd/go-toolbox/graceful"
	"github.com/egnd/go-toolbox/logging"
	"github.com/go-logr/zerologr"
)

var (
	version         = "debug"
	gracefulSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

var (
	showVersion = flag.Bool("version", false, "Show app version.")
	cfgPath     = flag.String("config", "configs/go-srv.yml", "Configuration file path.")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	zerologr.VerbosityFieldName = ""
	zerologr.NameFieldName = "service"

	ctx, cancel := graceful.Init(context.Background(), gracefulSignals...)
	defer cancel()

	cfg := config.NewViperCfg(config.ViperParams{Path: *cfgPath,
		UseOverride: true, OverrideSuffix: "override",
		UseEnv: true, EnvPrefix: "go-srv",
	})
	logger := logging.NewZerolog(logging.NewZerologCfgViper(cfg.Sub("logs")), os.Stderr)
	httpServer := services.NewGoFiber(ctx,
		cfg.Sub("server"), zerologr.New(&logger).WithName("httpserv"),
	)

	graceful.Register(httpServer.Start, httpServer.Stop)

	if err := graceful.Wait(); err != nil {
		logger.Fatal().Err(err).Msg("go-srv crashed")
	}

	logger.Info().Msg("go-srv stopped")
}
