package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
)

func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named:
		return []bool{isErrorType(t)}
	case *types.Pointer:
		return []bool{isErrorType(t)}
	case *types.Tuple:
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}

var errorType = types.
	Universe.Lookup("error").
	Type().
	Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	// проверяем, что t реализует интерфейс, при помощи которого определен тип error,
	// т.е. для типа t определен метод Error() string
	return types.Implements(t, errorType)
}

type Pass struct {
	Fset         *token.FileSet
	Files        []*ast.File
	OtherFiles   []string
	IgnoredFiles []string
	Pkg          *types.Package
	TypesInfo    *types.Info
}

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

var SuspiciousComparisonAnalyzer = &analysis.Analyzer{
	Name: "suscomp",
	Doc:  "detect suspicious comparisons like x == x",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if binExpr, ok := node.(*ast.BinaryExpr); ok {
					if binExpr.Op == token.EQL || binExpr.Op == token.NEQ {
						if pass.TypesInfo.Types[binExpr.X].Type == pass.TypesInfo.Types[binExpr.Y].Type &&
							pass.TypesInfo.Types[binExpr.X].Value == pass.TypesInfo.Types[binExpr.Y].Value {
							pass.Reportf(binExpr.Pos(), "suspicious self-comparison")
						}
					}
				}
				return true
			})
		}
		return nil, nil
	},
}

var InfiniteLoopAnalyzer = &analysis.Analyzer{
	Name: "infinite",
	Doc:  "check for potential infinite loops",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if forStmt, ok := node.(*ast.ForStmt); ok {
					if forStmt.Cond == nil {
						pass.Reportf(forStmt.Pos(), "potential infinite loop")
					}
				}
				return true
			})
		}
		return nil, nil
	},
}

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalysis",
	Doc:  "check for direct os.Exit calls in the main function of the main package",
	Run:  runOsExitAnalysis,
}

func runOsExitAnalysis(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if fn, ok := node.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				ast.Inspect(fn.Body, func(bodyNode ast.Node) bool {
					if callExpr, ok := bodyNode.(*ast.CallExpr); ok {

						if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "os" && sel.Sel.Name == "Exit" {
								pass.Reportf(callExpr.Pos(), "direct call to os.Exit in main function is prohibited")
							}
						}
					}
					return true
				})
			}
			return true
		})
	}
	return nil, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isReturnError(pass, call) {
				pass.Reportf(x.Pos(), "expression returns unchecked error")
			}
		}
	}
	tuplefunc := func(x *ast.AssignStmt) {
		if call, ok := x.Rhs[0].(*ast.CallExpr); ok {
			results := resultErrors(pass, call)
			for i := 0; i < len(x.Lhs); i++ {
				if id, ok := x.Lhs[i].(*ast.Ident); ok && id.Name == "_" && results[i] {
					pass.Reportf(id.NamePos, "assignment with unchecked error")
				}
			}
		}
	}
	errfunc := func(x *ast.AssignStmt) {
		for i := 0; i < len(x.Lhs); i++ {
			if id, ok := x.Lhs[i].(*ast.Ident); ok {
				if call, ok := x.Rhs[i].(*ast.CallExpr); ok {
					if id.Name == "_" && isReturnError(pass, call) {
						pass.Reportf(id.NamePos, "assignment with unchecked error")
					}
				}
			}
		}
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.ExprStmt:
				expr(x)
			case *ast.AssignStmt:

				if len(x.Rhs) == 1 {
					tuplefunc(x)
				} else {
					errfunc(x)
				}
			}
			return true
		})
	}
	return nil, nil
}
func main() {
	passes := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		shadow.Analyzer,
		sigchanyzer.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		unmarshal.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	allAnalyzers := append(passes,
		ErrCheckAnalyzer,
		InfiniteLoopAnalyzer,
		SuspiciousComparisonAnalyzer,
		printf.Analyzer,
		OsExitAnalyzer,
	)

	multichecker.Main(allAnalyzers...)
}
