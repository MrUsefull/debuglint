package logidentify_test

import (
	"debuglint/internal/configs"
	"debuglint/internal/funcs"
	"debuglint/pkg/logidentify"
	"path"
	"testing"

	"github.com/hashicorp/go-set/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAll verifies the overall analyzer against everything in the testdata directory.
func TestAll(t *testing.T) {
	tests := []struct {
		name         string
		testdataPath string
		want         *funcs.DebugFuncs
		wantErr      error
	}{
		{
			name:         "Finds all zap wrappers",
			testdataPath: path.Join(analysistest.TestData(), "zaplogger/zapwrappers"),
			want: func() *funcs.DebugFuncs {
				lf := funcs.NewDebugFuncs()
				lf.AddRequiresGuard(funcs.Description{
					Package: "zapwrappers",
					Name:    "ComplexSugaredRangeWrapperFunc",
				})
				lf.AddRequiresGuard(funcs.Description{
					Package: "zapwrappers",
					Name:    "ComplexSugaredForWrapperFunc",
				})
				lf.AddNormal(funcs.Description{
					Package: "zapwrappers",
					Name:    "SimpleSugaredWrapperFunc",
				})
				for f := range configs.ZapLogFunctions().Items() {
					lf.AddNormal(f)
				}
				for f := range configs.SlogLogFunctions().Items() {
					lf.AddNormal(f)
				}
				return lf
			}(),
		},
		{
			name:         "Finds all slog wrappers",
			testdataPath: path.Join(analysistest.TestData(), "sloglogger/slogwrappers"),
			want: func() *funcs.DebugFuncs {
				lf := funcs.NewDebugFuncs()
				lf.AddRequiresGuard(funcs.Description{
					Package: "slogwrappers",
					Name:    "ComplexRangeWrapper",
				})
				lf.AddRequiresGuard(funcs.Description{
					Package: "slogwrappers",
					Name:    "ComplexForWrapper",
				})
				lf.AddNormal(funcs.Description{
					Package: "slogwrappers",
					Name:    "SimpleWrapper",
				})
				lf.AddNormal(funcs.Description{
					Package: "slogwrappers",
					Name:    "SimpleWrapperCtx",
				})
				lf.AddNormal(funcs.Description{
					Package: "slogwrappers",
					Name:    "SlogObjDebugCtx",
				})
				lf.AddNormal(funcs.Description{
					Package: "slogwrappers",
					Name:    "SlogObjDebug",
				})
				for f := range configs.ZapLogFunctions().Items() {
					lf.AddNormal(f)
				}
				for f := range configs.SlogLogFunctions().Items() {
					lf.AddNormal(f)
				}
				return lf
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := analysistest.Run(t, tt.testdataPath, logidentify.Analyzer)
			require.Len(t, results, 1)

			returnVal, ok := results[0].Result.(*funcs.DebugFuncs)
			if assert.True(t, ok) {
				assertSetEquality(t, tt.want.All(), returnVal.All())
			}
		})
	}
}

// assertSetEquality gives more helpful output than just using assert.Equal
// against a funcs.DebugFuncs. Useful for finding typos in test cases.
func assertSetEquality[T comparable](tb testing.TB, want *set.Set[T], got *set.Set[T]) {
	tb.Helper()
	assert.Equal(tb, want.Size(), got.Size())

	if !assertContainsAll(tb, want, got) {
		tb.Logf("In want but not in got: %v\n", want.Difference(got))
		tb.Logf("In got but not in want: %v\n", got.Difference(want))
	}
}

func assertContainsAll[T comparable](tb testing.TB, want *set.Set[T], got *set.Set[T]) bool {
	tb.Helper()

	for found := range got.Items() {
		if !assert.True(tb, want.Contains(found), "%v not found in want: %v\n", found, want) {
			return false
		}
	}

	return true
}
