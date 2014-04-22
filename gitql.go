package main

import (
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/runtime"
	"github.com/cloudson/gitql/semantical"
	"path/filepath"
)

func main() {
	folder, errFile := filepath.Abs(*path)

	if errFile != nil {
		panic(errFile)
	}

	parser.New(query)
	ast, errGit := parser.AST()
	if errGit != nil {
		panic(errGit)
	}
	ast.Path = &folder
	errGit = semantical.Analysis(ast)
	if errGit != nil {
		panic(errGit)
	}

	runtime.Run(ast)
}
