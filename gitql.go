package main

import (
    "github.com/libgit2/git2go"
    "github.com/cloudson/gitql/parser"
    "github.com/cloudson/gitql/semantical"
    "github.com/cloudson/gitql/runtime"
    "fmt"
    "flag"
    "path/filepath"
)

func main() {
    path, errFile := filepath.Abs("./")
    if errFile != nil {
        panic(errFile)
    }
    query := flag.String("q", "", "Query ")
    flag.Parse()

    parser.New(*query)
    ast, errGit := parser.AST()
    if errGit != nil {
        panic(errGit)
    }
    ast.Path = query
    errGit = semantical.Analysis(ast)
    if errGit != nil {
        panic(errGit)
    }

    runtime.Run(ast)
    
    repo, err := git.OpenRepository(path)
    if err != nil {
        panic(err)
    }

    b, _ := git.NewOid("35bd9595f1d9a0a48f5c52fb923f0f2180f22976")
    obj, err2 := repo.LookupCommit(b)
    if err2 != nil {
        panic(err2)
    }
    fmt.Printf("\n%s\n",obj.Author())
}