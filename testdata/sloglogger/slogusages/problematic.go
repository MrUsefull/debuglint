package slogusages

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

func UnwrappedProblemUsage() {
	slog.Debug("some message", slog.String("important field", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	slog.Debug(expensiveStrFunc())                                                 // want `DL-1: Function call in unprotected debug log statement`
	slog.Debug("this is fine", slog.Int("intVal", 1), slog.String("stringVal", "stringKey"))
	for i := range 100 {
		// Do a thing..
		fmt.Println(i + i)
	}
}

func UnwrappedProblemUsageCtx(ctx context.Context) {
	slog.DebugContext(ctx, "some message", slog.String("important field", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	slog.DebugContext(ctx, expensiveStrFunc())                                                 // want `DL-1: Function call in unprotected debug log statement`
	slog.DebugContext(ctx, "this is fine", slog.Int("intVal", 1), slog.String("stringVal", "stringKey"))
}

func LoggerObjectProblemUsage() {
	logger.Debug("some message", slog.String("important field", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	logger.Debug(expensiveStrFunc())                                                 // want `DL-1: Function call in unprotected debug log statement`
	logger.Debug("this is fine", slog.Int("intVal", 1), slog.String("stringVal", "stringKey"))
}

func LoggerObjectProblemUsageCtx(ctx context.Context) {
	logger.DebugContext(ctx, "some message", slog.String("important field", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	logger.DebugContext(ctx, expensiveStrFunc())                                                 // want `DL-1: Function call in unprotected debug log statement`
	logger.DebugContext(ctx, "this is fine", slog.Int("intVal", 1), slog.String("stringVal", "stringKey"))
}

func AddsAndLogsExpensiveObj(a int, b int) int {
	out := a + b
	logger.Debug("result", slog.Int("out", out), slog.String("expensive", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	return out
}

func AddsAndLogsExpensiveFn(a int, b int) int {
	out := a + b
	slog.Debug("result", slog.Int("out", out), slog.String("expensive", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	return out
}

func expensiveStrFunc() string {
	sb := strings.Builder{}
	for i := range 1000 {
		sb.WriteString(fmt.Sprintf("-expensive%v-", i))
	}

	return sb.String()
}
