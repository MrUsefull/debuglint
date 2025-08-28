package funcs_test

import (
	"debuglint/internal/funcs"
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescriptionFromCall(t *testing.T) {
	t.Parallel()

	// Create a mock types.Func
	testPkg := types.NewPackage("test/package", "testpkg")
	testFunc := types.NewFunc(token.NoPos, testPkg, "TestFunc", &types.Signature{})
	testSel := &ast.SelectorExpr{
		X:   &ast.Ident{Name: "obj"},
		Sel: &ast.Ident{Name: "TestFunc"},
	}

	tests := []struct {
		name     string
		callExpr *ast.CallExpr
		uses     map[*ast.Ident]types.Object
		want     funcs.Description
		wantOk   bool
	}{
		{
			name: "valid selector expression with func in uses map",
			callExpr: &ast.CallExpr{
				Fun: testSel,
			},
			uses: map[*ast.Ident]types.Object{
				testSel.Sel: testFunc,
			},
			want: funcs.Description{
				Package: "test/package",
				Name:    "test/package.TestFunc",
			},
			wantOk: true,
		},
		{
			name: "non-selector expression",
			callExpr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "simpleFunc"},
			},
			uses: map[*ast.Ident]types.Object{},
		},
		{
			name: "selector not found in uses map",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "obj"},
					Sel: &ast.Ident{Name: "UnknownFunc"},
				},
			},
			uses: map[*ast.Ident]types.Object{},
		},
		{
			name: "selector maps to non-func object",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "obj"},
					Sel: &ast.Ident{Name: "NotAFunc"},
				},
			},
			uses: map[*ast.Ident]types.Object{
				{Name: "NotAFunc"}: types.NewVar(token.NoPos, testPkg, "NotAFunc", types.Typ[types.Int]),
			},
		},
		{
			name: "nil input call",
			uses: map[*ast.Ident]types.Object{
				{Name: "NotAFunc"}: types.NewVar(token.NoPos, testPkg, "NotAFunc", types.Typ[types.Int]),
			},
		},
		{
			name: "nil uses map",
			callExpr: &ast.CallExpr{
				Fun: testSel,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := funcs.DescriptionFromCall(tt.callExpr, tt.uses)

			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}
