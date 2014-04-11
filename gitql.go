package main

import (
    _"github.com/libgit2/git2go"
    "github.com/cloudson/gitql/parser"
    "github.com/cloudson/gitql/semantical"
    "github.com/cloudson/gitql/runtime"
    _"fmt"
    "flag"
    "path/filepath"
)

func main() {
    
    query := flag.String("q", "select * from commits", "The Query to search")
    pathString := flag.String("p", ".", "The (optional) path to run gitql")
    flag.Parse()

    path, errFile := filepath.Abs(*pathString)
    
    if errFile != nil {
        panic(errFile)
    }

    parser.New(*query)
    ast, errGit := parser.AST()
    if errGit != nil {
        panic(errGit)
    }
    ast.Path = &path
    errGit = semantical.Analysis(ast)
    if errGit != nil {
        panic(errGit)
    }

    runtime.Run(ast)
}
