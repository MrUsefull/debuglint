package configs

import (
	"debuglint/internal/funcs"

	"github.com/hashicorp/go-set/v3"
)

const (
	slogPackagePath = "log/slog"
	slogFieldType   = "log/slog.Attr"
)

// SlogLogFunctions provides the default set of known slog
// debug log functions.
func SlogLogFunctions() *set.Set[funcs.Description] {
	return set.From([]funcs.Description{
		{
			Package: slogPackagePath,
			Name:    "log/slog.Debug",
		},
		{
			Package: slogPackagePath,
			Name:    "log/slog.DebugContext",
		},
		{
			Package: slogPackagePath,
			Name:    "(*log/slog.Logger).Debug",
		},
		{
			Package: slogPackagePath,
			Name:    "(*log/slog.Logger).DebugContext",
		},
	})
}

// AllowedSlogCalls is the default set of allowed slog
// structured field calls.
func AllowedSlogCalls() *set.Set[funcs.Description] {
	return findStructuredLogFns(slogPackagePath, slogFieldType)
}
