package configs

import (
	"debuglint/internal/funcs"

	"github.com/hashicorp/go-set/v3"
)

// Config defines the full linter configuration.
type Config struct {
	// DebugLogFns describes logging functions that are
	// configured to be recognized by default.
	// Think zap's debug log functionality, or slog.Debug.
	DebugLogFns *set.Set[funcs.Description]

	// AllowedFnCalls is the set of function calls allowed
	// without a wrapping DebugEnabled check.
	AllowedFnCalls *set.Set[funcs.Description]

	// CustomDebugEnabledFns is a set of custom checks for
	// if debug logging is enabled. This is useful for
	// codebases that wrap logging libraries.
	CustomDebugEnabledFns *set.Set[funcs.Description]
}

// DefaultConfig returns the default configuration for linter.
func DefaultConfig() Config {
	return Config{
		DebugLogFns:           defaultKnownDebugLogFns(),
		AllowedFnCalls:        defaultAllowedFnCalls(),
		CustomDebugEnabledFns: &set.Set[funcs.Description]{},
	}
}

func defaultKnownDebugLogFns() *set.Set[funcs.Description] {
	out := ZapLogFunctions()
	out.InsertSet(SlogLogFunctions())

	return out
}

func defaultAllowedFnCalls() *set.Set[funcs.Description] {
	out := AllowedZapFnCalls()
	out.InsertSet(AllowedSlogCalls())

	return out
}
