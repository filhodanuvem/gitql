package runtime

import (
    "fmt"
    "github.com/libgit2/git2go"
)

var repo *git.Repository
type GitBuilder struct {
    tables map[string]string 
    possibleTables map[string][]string
}

type RuntimeError struct {
    code uint8 
    message string
}

func (e *RuntimeError) Error() string{
    return e.message
}

func throwRuntimeError(message string, code uint8) (*RuntimeError) {
    e := new(RuntimeError)
    e.message = message
    e.code = code

    return e
}

func GetGitBuilder(path string) (*GitBuilder) {
    gb := new(GitBuilder)
    gb.tables = make(map[string]string)
    possibleTables := map[string][]string {
        "commits": {
            "hash",
            "date",
            "author",
            "commiter",
            "message",
            "full_message",
        }, 
        "author": {
            "name",
            "email",
        },
        "files": {
            "hash",
            "path",
        },
    }
    gb.possibleTables = possibleTables

    return gb
}



func openRepository(path string) {
    _repo, err := git.OpenRepository(path)
    if err != nil {
        panic(err)
    }
    repo = _repo
}

func (g *GitBuilder) WithTable(tableName string, alias string) error {
    err := g.isValidTable(tableName)
    if err != nil {
        return err
    }

    if g.possibleTables[tableName] == nil {
        return throwRuntimeError(fmt.Sprintf("Table '%s' not found", tableName), 0)
    }

    if alias == "" {
        alias = tableName
    }

    g.tables[alias] = tableName 

    return nil
}

func (g *GitBuilder) isValidTable(tableName string) error {
    if g.possibleTables[tableName] == nil {
        return throwRuntimeError(fmt.Sprintf("Table '%s' not found", tableName), 0)
    }

    return nil
}

func (g *GitBuilder) UseFieldFromTable(field string, tableName string) error {
    err := g.isValidTable(tableName)
    if err != nil {
        return err
    }

    table := g.possibleTables[tableName]
    for _, t := range table {
        if t == tableName {
            return nil
        }
    }

    return throwRuntimeError(fmt.Sprintf("Table '%s' has not field '%s'", tableName, field), 0)
}

