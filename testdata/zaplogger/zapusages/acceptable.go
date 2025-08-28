package zapusages

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// NormalDebugCheck should not trigger a debug log warning since it does other stuff
// and it checks for debug enabled via check
func NormalDebugCheck(in []int) {
	for i := range in {
		fmt.Println(i)
		if ce := logger.Check(zap.DebugLevel, "hello"); ce != nil {
			ce.Write(zap.Int("i", i))
		}
	}
}

// NormalDebugLevelCheck should not trigger a debug log warning since it does other stuff
// and it checks for debug enabled by comparing debug level
func NormalDebugLevelCheck(in []int) {
	for i := range in {
		fmt.Println(i)
		if logger.Level() <= zapcore.DebugLevel {
			logger.Debug("hello", zap.Int("i", i))
		}
	}
}

func AddsAndLogs(a int, b int) int {
	out := a + b
	logger.Debug("result", zap.Int("out", out))
	return out
}
