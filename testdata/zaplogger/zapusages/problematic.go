package zapusages

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func UnwrappedProblemUsage() {
	logger.Debug("some message", zap.String("important field", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	logger.Debug(expensiveStrFunc())                                                // want `DL-1: Function call in unprotected debug log statement`
	logger.Debug("this is fine", zap.Int("intVal", 1), zap.String("stringVal", "stringKey"))
	for i := range 100 {
		// Do a thing..
		fmt.Println(i + i)
	}
}

func AddsAndLogsExpensive(a int, b int) int {
	out := a + b
	logger.Debug("result", zap.Int("out", out), zap.String("expensive", expensiveStrFunc())) // want `DL-1: Function call in unprotected debug log statement`
	return out
}

func expensiveStrFunc() string {
	sb := strings.Builder{}
	for i := range 1000 {
		sb.WriteString(fmt.Sprintf("-expensive%v-", i))
	}

	return sb.String()
}
