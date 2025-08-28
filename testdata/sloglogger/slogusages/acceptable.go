package slogusages

import (
	"context"
	"fmt"
	"log/slog"
)

// NormalDebugEnabled should not trigger a debug log warning since it does other stuff
// and it checks for debug enabled via Enabled check
func NormalDebugEnabled(ctx context.Context, in []int) {
	for i := range in {
		fmt.Println(i)
		if logger.Enabled(ctx, slog.LevelDebug) {
			logger.Debug("hello", slog.Int("i", i))
		}
	}
}

// NormalDebugEnabledSlogDefault should not trigger a debug log warning since it does other stuff
// and it checks for debug enabled using default slog
func NormalDebugEnabledSlogDefault(ctx context.Context, in []int) {
	for i := range in {
		fmt.Println(i)
		if slog.Default().Enabled(ctx, slog.LevelDebug) {
			slog.Debug("hello", slog.Int("i", i))
		}
	}
}

func NormalDebugLevelCheck(ctx context.Context, in []int) {
	logger.DebugContext(ctx, "hello", slog.Any("in", in))
}

func AddsAndLogsFn(a int, b int) int {
	out := a + b
	slog.Debug("result", slog.Int("out", out))
	return out
}

func AddsAndLogsObj(a int, b int) int {
	out := a + b
	logger.Debug("result", slog.Int("out", out))
	return out
}
