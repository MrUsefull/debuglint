package configs

import (
	"debuglint/internal/funcs"

	"github.com/hashicorp/go-set/v3"
)

const (
	zapPackagePath = "go.uber.org/zap"
	zapFieldType   = "go.uber.org/zap.Field"
)

// ZapLogFunctions returns the known set of zap log functions at compile time.
func ZapLogFunctions() *set.Set[funcs.Description] {
	return set.From([]funcs.Description{
		{
			Package: zapPackagePath,
			Name:    "(*go.uber.org/zap.SugaredLogger).Debug",
		},
		{
			Package: zapPackagePath,
			Name:    "(*go.uber.org/zap.Logger).Debug",
		},
	})
}

// AllowedZapFnCalls returns the default zap function calls that are allowed.
func AllowedZapFnCalls() *set.Set[funcs.Description] {
	return findStructuredLogFns(zapPackagePath, zapFieldType)
}
