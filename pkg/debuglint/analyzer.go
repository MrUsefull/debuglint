// Package debuglint provides an analyzer that lints debug statements.
package debuglint

import (
	"debuglint/internal/configs"
	"debuglint/internal/funcs"
	"debuglint/pkg/linterrs"
	"debuglint/pkg/logidentify"
	"fmt"
	"go/ast"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// DL1 is the lint message when expensive function calls are made without
// any log guard.
const DL1 string = "DL-1: Function call in unprotected debug log statement"

// Analyzer is the DebugLint Analyzer.
//
//nolint:gochecknoglobals // Analyzer must be a package global
var Analyzer = &analysis.Analyzer{
	Name:     "DebugLint",
	Doc:      "Identifies performance impacting debug log statements",
	URL:      "",
	Run:      run,
	Requires: []*analysis.Analyzer{logidentify.Analyzer, inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	knownLogFns, ok := pass.ResultOf[logidentify.Analyzer].(*funcs.DebugFuncs)
	if !ok {
		return nil, fmt.Errorf("logidentify: %w", linterrs.ErrRequirementsFailed)
	}

	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("inspector: %w", linterrs.ErrRequirementsFailed)
	}

	analyzer := newAnalyzer(pass, knownLogFns, configs.DefaultConfig())

	inspector.Preorder([]ast.Node{
		(*ast.FuncDecl)(nil),
	}, analyzer.inspectFuncDecl)

	//nolint:nilnil // this is the correct way
	return nil, nil
}

type analyzer struct {
	pass          *analysis.Pass
	knownDebugFns *funcs.DebugFuncs
	cfg           configs.Config
}

func newAnalyzer(
	pass *analysis.Pass,
	knownDebugFns *funcs.DebugFuncs,
	cfg configs.Config,
) *analyzer {
	return &analyzer{
		pass:          pass,
		knownDebugFns: knownDebugFns,
		cfg:           cfg,
	}
}

func (a *analyzer) inspectFuncDecl(n ast.Node) {
	fnDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		// should never occur
		return
	}

	a.inspectStmts(fnDecl.Body.List)
}

// inspectStmts recursively examines statements looking for unguarded debug function calls.
// It only processes expression statements at the top level to avoid false positives
// within control flow structures that may have their own guards.
func (a *analyzer) inspectStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		switch x := stmt.(type) {
		case *ast.ExprStmt:
			a.inspectUnguardedExpr(x.X)
		default:
		}
	}
}

func (a *analyzer) inspectUnguardedExpr(n ast.Expr) {
	fnCall, ok := n.(*ast.CallExpr)
	if !ok {
		return
	}

	fn, ok := funcs.DescriptionFromCall(fnCall, a.pass.TypesInfo.Uses)
	if !ok {
		return
	}

	if !a.knownDebugFns.All().Contains(fn) {
		// not a debug fn call, no need to continue
		return
	}

	// At this point we are sure this is a debug function call
	// that is not guarded. Inspect our args for disallowed
	// function calls.
	if a.argsHaveFnCall(fnCall.Args) {
		a.pass.Report(analysis.Diagnostic{
			Pos:     fnCall.Pos(),
			End:     fnCall.End(),
			Message: DL1,
		})
	}
}

func (a *analyzer) argsHaveFnCall(args []ast.Expr) bool {
	return slices.ContainsFunc(args, a.inspectOneArg)
}

// inspectOneArg recursively analyzes a single function call argument for violations.
// Returns true if the argument contains an expensive function call that should be guarded.
// Allowed functions (like zap.String) are permitted, but their arguments are still validated.
func (a *analyzer) inspectOneArg(arg ast.Expr) bool {
	fnCallArg, ok := arg.(*ast.CallExpr)
	if !ok {
		return false
	}

	fnDesc, ok := funcs.DescriptionFromCall(fnCallArg, a.pass.TypesInfo.Uses)
	if !ok {
		return true
	}

	if !a.cfg.AllowedFnCalls.Contains(fnDesc) {
		return true
	}

	return slices.ContainsFunc(fnCallArg.Args, a.inspectOneArg)
}
