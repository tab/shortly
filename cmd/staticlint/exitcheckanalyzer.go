package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ExitCheckAnalyzer checks for direct os.Exit calls
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for direct os.Exit calls",
	Run:  run,
}

// run checks for direct os.Exit calls
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkgIdent, ok := selExpr.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "os" || selExpr.Sel.Name != "Exit" {
				return true
			}

			pos := pass.Fset.Position(callExpr.Pos())
			pass.Reportf(callExpr.Pos(), "direct call to os.Exit found in %s:%d", pos.Filename, pos.Line)
			return false
		})
	}

	return nil, nil
}
