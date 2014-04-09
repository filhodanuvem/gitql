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
    path, errFile := filepath.Abs("./")
    
    if errFile != nil {
        panic(errFile)
    }
    query := flag.String("q", "select * from commits", "The Query to search")
    flag.Parse()

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
