//go:build tools

// tools package exists to keep certain build dependencies
// in the go.mod files after running "go mod tidy"
package tools

import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "go.uber.org/zap"
)
