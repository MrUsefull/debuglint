package configs

import (
	"debuglint/internal/funcs"

	"github.com/MrUsefull/FuncFind/pkg/funcfind"
	"github.com/hashicorp/go-set/v3"
)

func findStructuredLogFns(pkgPath string, returnType string) *set.Set[funcs.Description] {
	fnIter, err := funcfind.Returning(pkgPath, returnType)
	if err != nil {
		panic(err)
	}

	s := set.New[funcs.Description](1)
	for fn := range fnIter {
		s.Insert(
			funcs.Description{
				Package: pkgPath,
				Name:    fn.FullName(),
			},
		)
	}

	return s
}
