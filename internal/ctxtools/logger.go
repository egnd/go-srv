package ctxtools

import (
	"context"

	"github.com/go-logr/logr"
)

func AddLogger(ctx context.Context, logger logr.Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

func GetLogger(ctx context.Context) logr.Logger {
	return logr.FromContextOrDiscard(ctx)
}
