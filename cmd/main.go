// package main runs the linter
package main

import (
	"debuglint/pkg/debuglint"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(debuglint.Analyzer)
}
