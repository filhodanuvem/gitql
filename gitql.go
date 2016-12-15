package main

import (
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/runtime"
	"github.com/cloudson/gitql/semantical"
	"log"
	"path/filepath"
)

func main() {
	folder, errFile := filepath.Abs(*path)

	if errFile != nil {
		log.Fatalln(errFile)
	}

	parser.New(query)
	ast, errGit := parser.AST()
	if errGit != nil {
		log.Fatalln(errGit)
	}
	ast.Path = &folder
	errGit = semantical.Analysis(ast)
	if errGit != nil {
		log.Fatalln(errGit)
	}

	runtime.Run(ast, genJson)
}
