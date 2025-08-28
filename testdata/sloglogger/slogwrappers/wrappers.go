package slogwrappers

import (
	"context"
	"log/slog"
	"os"
)

func SimpleWrapper(someStr string) {
	slog.Debug(someStr, slog.String("key", "value"))
}

func SimpleWrapperCtx(ctx context.Context, someStr string) {
	slog.DebugContext(ctx, someStr, slog.String("key", "value"))
}

func ComplexRangeWrapper(things []string) {
	for _, thing := range things {
		slog.Debug("msg", slog.String("thing", thing))
	}
}

func ComplexForWrapper(things []string) {
	for i := 0; i < len(things); i++ {
		slog.Debug("msg", slog.String("thing", things[i]))
	}
}

var logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

func SlogObjDebug(fieldStr string) {
	logger.Debug("some msg", slog.String("key", fieldStr))
}

func SlogObjDebugCtx(ctx context.Context, fieldStr string) {
	logger.DebugContext(ctx, "some msg", slog.String("key", fieldStr))
}
