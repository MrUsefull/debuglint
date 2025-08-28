package debuglint_test

import (
	"debuglint/pkg/debuglint"
	"path"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name         string
		testdataPath string
		wantErr      error
	}{
		{
			name:         "Test zap usage",
			testdataPath: path.Join(analysistest.TestData(), "zaplogger/zapusages"),
		},
		{
			name:         "Test slog usage",
			testdataPath: path.Join(analysistest.TestData(), "sloglogger/slogusages"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysistest.Run(t, tt.testdataPath, debuglint.Analyzer)
		})
	}
}
