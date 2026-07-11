package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "chechosexitinmain",
	Doc:  "запрещает прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	var mainFunc *ast.FuncDecl
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" && x.Recv == nil {
					mainFunc = x
				}
			}
			return true
		})
	}
	if mainFunc == nil {
		return nil, nil
	}
	ast.Inspect(mainFunc.Body, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if s, ok := c.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := s.X.(*ast.Ident); ok {
					if ident.Name == "os" && s.Sel.Name == "Exit" {
						pass.Reportf(c.Pos(), "прямой вызов os.Exit в функции main запрещен")
					}
				}
			}
		}
		return true
	})
	return nil, nil
}
