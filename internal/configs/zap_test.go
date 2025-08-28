package configs

import (
	"debuglint/internal/funcs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllowedZapFnCalls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		wantInclude []funcs.Description
	}{
		{
			name: "Spot check functions",
			wantInclude: []funcs.Description{
				{
					Package: zapPackagePath,
					Name:    "go.uber.org/zap.Int",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AllowedZapFnCalls()
			assert.True(t, got.ContainsSlice(tt.wantInclude), "Got: %v\nWant:%v\n", got, tt.wantInclude)
		})
	}
}
