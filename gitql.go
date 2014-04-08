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
    path, errFile := filepath.Abs("/home/cloud/Rocket/braprint")
    
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
    ast.Path = &path
    errGit = semantical.Analysis(ast)
    if errGit != nil {
        panic(errGit)
    }

    runtime.Run(ast)
    


    // ================ Testing 
    // repo, err := git.OpenRepository(path)
    // if err != nil {
    //     panic(err)
    // }

    // b, _ := git.NewOid("35bd9595f1d9a0a48f5c52fb923f0f2180f22976")
    // obj, err2 := repo.LookupCommit(b)
    // if err2 != nil {
    //     panic(err2)
    // }
    // fmt.Printf("\n%s\n",obj.Author())

    // walk, err:= repo.Walk()
    // if err != nil {
    //     fmt.Printf(err.Error())
    // }
    // walk.PushHead()
    // walk.Sorting(git.SortTime)
    // i := 1
    // fn := func (oid *git.Commit) bool {
    //     fmt.Printf("\n%s\n", oid.Id())
    //     i = i + 1
    //     if i == 5 {
    //         return false
    //     }
    //     return true
    // }
    // err = walk.Iterate(fn)
    // if err != nil {
    //     fmt.Printf(err.Error())
    // }

}