// Package logfuncs contains types used to identify log functions
package funcs_test

import (
	"debuglint/internal/funcs"
	"testing"

	"github.com/hashicorp/go-set/v3"
	"github.com/stretchr/testify/assert"
)

func TestDebugFuncs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		optNormal   []funcs.Description
		addNormal   []funcs.Description
		addRequired []funcs.Description
	}{
		{
			name: "Basic usage",
			optNormal: []funcs.Description{
				{
					Package: "go.uber.org/zap",
					Name:    "go.uber.org/zap.Int",
				},
				{
					Package: "go.uber.org/zap",
					Name:    "go.uber.org/zap.String",
				},
			},
			addNormal: []funcs.Description{
				{
					Package: "my.normal.package/here",
					Name:    "MyNormalFunctionNameHere",
				},
			},
			addRequired: []funcs.Description{
				{
					Package: "my.required.package/here",
					Name:    "MyRequiredFunctionNameHere",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lf := funcs.NewDebugFuncs(funcs.WithNormalDebugFns(set.From(tt.optNormal)))
			for _, f := range tt.addNormal {
				lf.AddNormal(f)
			}

			for _, f := range tt.addRequired {
				lf.AddRequiresGuard(f)
			}

			assert.ElementsMatch(t, lf.All().Slice(), append(append(tt.optNormal, tt.addNormal...), tt.addRequired...))
			assert.ElementsMatch(t, lf.Normal().Slice(), append(tt.optNormal, tt.addNormal...))
			assert.ElementsMatch(t, lf.RequiresGuard().Slice(), tt.addRequired)
		})
	}
}
