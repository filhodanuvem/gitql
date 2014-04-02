package runtime

import (
    "fmt"
    "reflect"
    "github.com/cloudson/gitql/parser"
    "github.com/libgit2/git2go"
)

var repo *git.Repository
var builder *GitBuilder

type GitBuilder struct {
    tables map[string]string 
    possibleTables map[string][]string
}

type RuntimeError struct {
    code uint8 
    message string
}

type RuntimeVisitor struct {

}

// =========================== Runtime
func Run(n *parser.NodeProgram) {
    builder = GetGitBuilder(n.Path)
    visitor := new(RuntimeVisitor)
    err := visitor.Visit(n)
    if err != nil {
        panic(err)
    }

    // builder := visitor.Builder()
}

// =========================== Error 

func (e *RuntimeError) Error() string{
    return e.message
}

func throwRuntimeError(message string, code uint8) (*RuntimeError) {
    e := new(RuntimeError)
    e.message = message
    e.code = code

    return e
}

// ========================== RuntimeVisitor

func (v *RuntimeVisitor) Visit(n *parser.NodeProgram) (error) {
    return v.VisitSelect(n.Child.(*parser.NodeSelect))
} 

func (v *RuntimeVisitor) VisitSelect(n *parser.NodeSelect) (error) {

    return nil 
} 

func (v *RuntimeVisitor) VisitExpr(n parser.NodeExpr) (error) {
    switch reflect.TypeOf(n) {
        case reflect.TypeOf(new(parser.NodeGreater)) : 
            g:= n.(*parser.NodeGreater)
            return v.VisitGreater(g)
        case reflect.TypeOf(new(parser.NodeSmaller)) : 
            g:= n.(*parser.NodeSmaller)
            return v.VisitSmaller(g)
    } 

    return nil
}

func (v *RuntimeVisitor) VisitGreater(n *parser.NodeGreater) (error) {
    
    return nil
}

func (v *RuntimeVisitor) VisitSmaller(n *parser.NodeSmaller) (error) {
    
    return nil
}

func (v *RuntimeVisitor) Builder() (*GitBuilder){
    return nil
}


// =================== GitBuilder 

func GetGitBuilder(path *string) (*GitBuilder) {

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



func openRepository(path *string) {
    _repo, err := git.OpenRepository(*path)
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

