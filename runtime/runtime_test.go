package runtime

import (
    "testing"
)

func TestErrorWithInvalidTables(t *testing.T) {
    invalidTables := []string{
        "cloudson",
        "blah",
    }

    gb := GetGitBuilder("./")
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
        "author",
    }

    gb := GetGitBuilder("./")
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
        {"author", "blah"},
        {"files", ""},
    }

    gb := GetGitBuilder("./")
    for _, tableMetada := range metadata {
        err := gb.UseFieldFromTable(tableMetada[1], tableMetada[0])
        if err == nil {
            t.Errorf("Table '%s' should has not field '%s'", tableMetada[0], tableMetada[1])
        }
    }
}