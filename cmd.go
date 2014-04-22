package main

import (
    "flag"
    "os"
    "fmt"
)

var path *string
var query string

func init() {
    parseCommandLine()
}

func usage() {
    fmt.Println("Gitql - Git query language")
    fmt.Printf("Usage: %s [flags] [arg] \n ", os.Args[0])
    fmt.Printf("\nFlags: \n")
    flag.PrintDefaults()
    fmt.Printf("Arguments: \n")
    fmt.Printf("  path: Path directory of a git repository\n")
}

func parseCommandLine() {
    path = flag.String("p", ".", "The (optional) path to run gitql")
    version := flag.Bool("v", false, "The version of gitql")
    flag.Usage = usage
    flag.Parse()

    if *version {
        // @todo refactor to dynamic value
        fmt.Println("Gitql 1.0.0-RC4")
        os.Exit(0)
    }

    query = flag.Arg(0)
    if flag.NArg() != 1 {
        flag.Usage()
        os.Exit(1)
    }
}