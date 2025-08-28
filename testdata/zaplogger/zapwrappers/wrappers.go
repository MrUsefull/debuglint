// Package zapwrappers is used to test for zap log wrapper detection
package zapwrappers

import (
	"fmt"

	"go.uber.org/zap"
)

// globalSugarLogger is used  to verify we can detect direct usages of a zap logger variable
var globalSugarLogger = func() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("logger failed: %v", err))
	}
	return logger.Sugar()
}()

// SimpleSugaredWrapperFunc should be detected as a debug log function
// by the logidentify analyzer
func SimpleSugaredWrapperFunc(field string) {
	globalSugarLogger.Debug("some standard message", zap.String("field", field))
}

// ComplexSugaredRangeWrapperFunc is an example debug log function that ranges
// over a collection, and only debug logs
func ComplexSugaredRangeWrapperFunc(things []string) {
	for _, thing := range things {
		globalSugarLogger.Debug("some message", zap.String("thing", thing))
	}
}

// ComplexSugaredForWrapperFunc is an example debug log function that ranges
// over a collection, and only debug logs
func ComplexSugaredForWrapperFunc(things []string) {
	for i := 0; i < len(things); i++ { //nolint:rangelint // for loop is the point here
		globalSugarLogger.Debug("some message", zap.String("thing", things[i]))
	}
}
