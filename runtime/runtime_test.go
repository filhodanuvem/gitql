package runtime

import (
    "path/filepath"
    "testing"
)

func TestErrorWithInvalidTables(t *testing.T) {
    invalidTables := []string{
        "cloudson",
        "blah",
    }

    var path string 
    path, _ = filepath.Abs("../")
    gb := GetGitBuilder(&path)
    for _, tableName := range invalidTables {
        err := gb.WithTable(tableName, tableName)
        if err == nil {
            t.Errorf("Table '%s' should throws an error", tableName)
        }
    }
}

func TestTablesWithoutAlias(t *testing.T) {
    invalidTables := []string{
        "commits",
        "trees",
    }

    var path string 
    path, _ = filepath.Abs("../")
    gb := GetGitBuilder(&path)
    for _, tableName := range invalidTables {
        err := gb.WithTable(tableName, "")
        if err != nil {
            t.Errorf(err.Error())
        }
    }   
}

func TestNotFoundFieldsFromTable(t *testing.T) {
    metadata := [][]string {
        {"commits", "hashas"},
        {"trees", "blah"},
        {"commits", ""},
    }

    var path string 
    path, _ = filepath.Abs("../")
    gb := GetGitBuilder(&path)
    for _, tableMetada := range metadata {
        err := gb.UseFieldFromTable(tableMetada[1], tableMetada[0])
        if err == nil {
            t.Errorf("Table '%s' should has not field '%s'", tableMetada[0], tableMetada[1])
        }
    }
}

func TestFoundFieldsFromTable(t *testing.T) {
    metadata := [][]string {
        {"commits", "*"},
        {"trees", "hash"},
        {"commits", "hash"},
    }

    var path string 
    path, _ = filepath.Abs("../")
    gb := GetGitBuilder(&path)
    for _, tableMetada := range metadata {
        err := gb.UseFieldFromTable(tableMetada[1], tableMetada[0])
        if err != nil {
            t.Errorf(err.Error())
        }
    }
}
