package funcs

import (
	"go/ast"
	"go/types"
)

// Description contains a description of a Function.
// Description is used instead of *types.Func for
// ease of configuration.
type Description struct {
	// Package is the package the function is declared in.
	// Should be the full package path.
	Package string
	// Name is the name of the function.
	Name string
}

// DescriptionFromCall returns a FuncDescription built from the
// provided ast.CallExpr and uses map. Returns the found FuncDescription and true if
// the FuncDescription was able to be created.
func DescriptionFromCall(n *ast.CallExpr, uses map[*ast.Ident]types.Object) (Description, bool) {
	if n == nil || uses == nil {
		return Description{}, false
	}

	fn, ok := getFuncType(n, uses)
	if !ok {
		return Description{}, false
	}

	return descriptionFromFunc(fn), true
}

// descriptionFromFunc converts a *types.Func to a FuncDescription.
func descriptionFromFunc(fn *types.Func) Description {
	return Description{
		Name:    fn.FullName(),
		Package: fn.Pkg().Path(),
	}
}

func getFuncType(call *ast.CallExpr, uses map[*ast.Ident]types.Object) (*types.Func, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	obj, ok := uses[sel.Sel]
	if !ok {
		return nil, false
	}

	f, ok := obj.(*types.Func)

	return f, ok
}
