// Package logidentify provides an analyzer that identifies log functions
// and log function wrappers
package logidentify

import (
	"debuglint/internal/configs"
	"debuglint/internal/funcs"
	"debuglint/pkg/linterrs"
	"fmt"
	"go/ast"
	"reflect"

	"github.com/hashicorp/go-set/v3"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	// Analyzer is the logidentify analyzer.
	//
	//nolint:gochecknoglobals // Analyzer must be a package global
	Analyzer = &analysis.Analyzer{
		Name:       "LogIdentify",
		Doc:        "Identifies log functions",
		URL:        "",
		Run:        run,
		Requires:   []*analysis.Analyzer{inspect.Analyzer},
		ResultType: reflect.TypeOf((*funcs.DebugFuncs)(nil)),
	}
)

func run(pass *analysis.Pass) (any, error) {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("inspector: %w", linterrs.ErrRequirementsFailed)
	}

	analyzer := newAnalyzer(pass, configs.DefaultConfig().DebugLogFns)

	inspector.Preorder([]ast.Node{
		(*ast.FuncDecl)(nil),
	}, analyzer.inspectFuncDecl)

	return analyzer.foundLogFns, nil
}

type analyzer struct {
	pass *analysis.Pass

	foundLogFns *funcs.DebugFuncs

	currentFn funcs.Description

	knownDebugLogFns *set.Set[funcs.Description]
}

func newAnalyzer(pass *analysis.Pass, knownDebugLogFns *set.Set[funcs.Description]) *analyzer {
	return &analyzer{
		pass: pass,
		foundLogFns: funcs.NewDebugFuncs(
			funcs.WithNormalDebugFns(knownDebugLogFns),
		),
		knownDebugLogFns: knownDebugLogFns,
	}
}

func (a *analyzer) inspectFuncDecl(n ast.Node) {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return // should never occur
	}

	a.populateCurrentFuncDeclName(funcDecl)

	ops := &opCounts{}
	a.inspectFuncBody(funcDecl.Body, ops)

	if ops.hasOtherStmts || ops.numDebugStmts == 0 {
		return // not a debug function
	}

	if ops.numCtrlStmts > 0 {
		a.foundLogFns.AddRequiresGuard(a.currentFn)
	} else {
		a.foundLogFns.AddNormal(a.currentFn)
	}
}

// populateCurrentFuncDeclName extracts the function name and package from an AST node
// and stores it in the analyzer for use in function classification decisions.
func (a *analyzer) populateCurrentFuncDeclName(funcDecl *ast.FuncDecl) {
	fnDescription := funcs.Description{
		Name:    funcDecl.Name.Name,
		Package: a.pass.Pkg.Name(),
	}
	a.currentFn = fnDescription
}

// opCounts tracks the types of operations found in a function body during analysis.
// This helps categorize functions as debug wrappers vs regular functions.
type opCounts struct {
	// numDebugStmts is the count of debug logging statements.
	numDebugStmts int
	// numCtrlStmts is the count of control flow statements (for/range loops).
	numCtrlStmts int
	// hasOtherStmts is true if any non-debug, non-control statements found.
	hasOtherStmts bool
}

// inspectFuncBody recursively analyzes a function's statements to categorize it.
// Updates ops counts to determine if this is a debug wrapper function.
// Functions with only debug calls become "normal" debug functions.
// Functions with debug calls + control flow become "requiresGuard" debug functions.
// Functions with other statements are not considered debug functions.
func (a *analyzer) inspectFuncBody(n *ast.BlockStmt, ops *opCounts) {
	for _, stmt := range n.List {
		switch x := stmt.(type) {
		case *ast.ExprStmt:
			if a.isDebugCall(x) {
				ops.numDebugStmts++
			} else {
				ops.hasOtherStmts = true
			}
		case *ast.ForStmt:
			ops.numCtrlStmts++
			a.inspectFuncBody(x.Body, ops)
		case *ast.RangeStmt:
			ops.numCtrlStmts++
			a.inspectFuncBody(x.Body, ops)
		default:
			ops.hasOtherStmts = true
		}

		if ops.hasOtherStmts {
			// we found something that is not a debug call or a control statement
			return
		}
	}
}

func (a *analyzer) isDebugCall(n *ast.ExprStmt) bool {
	fnCall, ok := n.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	fn, ok := funcs.DescriptionFromCall(fnCall, a.pass.TypesInfo.Uses)
	if !ok {
		return false
	}

	return a.knownDebugLogFns.Contains(fn)
}
