package configs_test

import (
	"debuglint/internal/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	// Verify nothing nil. Additional cases
	// to be added as needed
	got := configs.DefaultConfig()
	assert.NotNil(t, got.AllowedFnCalls)
	assert.NotNil(t, got.CustomDebugEnabledFns)
	assert.NotNil(t, got.DebugLogFns)
}
