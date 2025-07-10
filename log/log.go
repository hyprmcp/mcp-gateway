package log

import (
	"context"
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type key struct{}

var loggerKey = key{}

var rootLogger = stdr.New(log.Default())

func Root() logr.Logger {
	return rootLogger
}

func Add(ctx context.Context, logger logr.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Get(ctx context.Context) logr.Logger {
	if ctx == nil {
		return rootLogger
	}

	logger, ok := ctx.Value(loggerKey).(logr.Logger)
	if !ok {
		return rootLogger
	}

	return logger
}
