package main

import (
    "flag"
    "os"
    "fmt"
    "github.com/cloudson/gitql/runtime"
)

var path *string
var query string
var genJson *bool

func init() {
    parseCommandLine()
}

func usage() {
    fmt.Println("Gitql - Git query language")
    fmt.Printf("Usage: %s [flags] [args] \n ", os.Args[0])
    fmt.Printf("\nFlags: \n")
    flag.PrintDefaults()
    fmt.Printf("Arguments: \n")
    fmt.Printf("  sql: A query to run\n")
}

func printTables() {
    tables := runtime.PossibleTables()
    for tableName, fields := range tables {
        fmt.Printf("%s\n\t", tableName)
        for i, field := range fields {
            comma := "."
            if i + 1 < len(fields) {
                comma = ", "
            }
            fmt.Printf("%s%s", field, comma)
        }
        fmt.Println()
    }

}

func parseCommandLine() {
    path = flag.String("p", ".", "The (optional) path to run gitql")
    version := flag.Bool("v", false, "The version of gitql")
    showTables := flag.Bool("show-tables", false, "Show all tables")
    genJson = flag.Bool("json", false, "Generate JSON")
    flag.Usage = usage
    flag.Parse()

    if *version {
        // @todo refactor to dynamic value
        fmt.Println("Gitql 1.1.1")
        os.Exit(0)
    }

    if *showTables {
        fmt.Printf("Tables: \n\n")
        printTables()
        os.Exit(0)
    }

    query = flag.Arg(0)
    if flag.NArg() != 1 {
        flag.Usage()
        os.Exit(1)
    }
}
