package runtime

import (
    "fmt"
    "reflect"
    "strings"
    "github.com/cloudson/gitql/parser"
    "github.com/libgit2/git2go"
)

var repo *git.Repository
var builder *GitBuilder
var boolRegister bool

type GitBuilder struct {
    tables map[string]string 
    possibleTables map[string][]string
    walk *git.RevWalk
    object *git.Commit
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

    s := n.Child.(*parser.NodeSelect)
    where := s.Where
    // valuesTable := make([]string, len(s.Fields))
    counter := 1
    fmt.Println()
    fn := func (object *git.Commit) bool {
        builder.setObject(object)
        visitor.VisitExpr(where)
        if boolRegister {
            fields := s.Fields
            if s.WildCard {
                fields = builder.possibleTables[s.Tables[0]]
            }            
            for _, f := range fields {
                fmt.Printf("%s | ", discoverLvalue(f, s.Tables[0], object))    
            }
            fmt.Println()

            
            counter = counter + 1
        }
        if counter > s.Limit {
            return false
        }
        return true
    }

    err = builder.walk.Iterate(fn)
    if err != nil {
        fmt.Printf(err.Error())
    }

}

func discoverLvalue(identifier string, table string, object *git.Commit) string {
    err := builder.UseFieldFromTable(identifier, table)
    if err != nil {
        panic(err)
    }
    switch identifier {
        case "hash" : 
            return object.Id().String()
        case "author":
            return object.Author().Name
        case "author_email":
            return object.Author().Email
        case "committer":
            return object.Committer().Name
        case "committer_email":    
            return object.Committer().Email
        case "date":
            return object.Committer().When.String()
        case "full_message":
            return object.Message()
        case "message": 
            message := object.Message()
            r := []rune("\n")
            idx := strings.IndexRune(message, r[0])
            if idx != -1 {
                message = message[0:idx]
            }
            return message  

    }

    panic(fmt.Sprintf("Trying select field %s ", identifier))
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
    table := n.Tables[0]
    fields := n.Fields 
    var err error
    for _, f := range fields {
        err = builder.UseFieldFromTable(f, table)
        if err != nil {
            return err
        }
    }
    // Why not visit expression right now ? 
    // Because we will, at first, discover the current commit 
    return nil 
} 

func (v *RuntimeVisitor) VisitExpr(n parser.NodeExpr) (error) {
    switch reflect.TypeOf(n) {
        case reflect.TypeOf(new(parser.NodeEqual)) : 
            g:= n.(*parser.NodeEqual)
            return v.VisitEqual(g)
        case reflect.TypeOf(new(parser.NodeGreater)) : 
            g:= n.(*parser.NodeGreater)
            return v.VisitGreater(g)
        case reflect.TypeOf(new(parser.NodeSmaller)) : 
            g:= n.(*parser.NodeSmaller)
            return v.VisitSmaller(g)
        case reflect.TypeOf(new(parser.NodeOr)) : 
            g:= n.(*parser.NodeOr)
            return v.VisitOr(g)
        case reflect.TypeOf(new(parser.NodeAnd)) : 
            g:= n.(*parser.NodeAnd)
            return v.VisitAnd(g)

    } 

    return nil
}

func (v *RuntimeVisitor) VisitEqual(n *parser.NodeEqual) (error) {
    lvalue := n.LeftValue().(*parser.NodeId).Value()
    rvalue := n.RightValue().(*parser.NodeLiteral).Value()
    boolRegister = n.Assertion(discoverLvalue(lvalue, "commits", builder.object), rvalue)
    
    return nil
}

func (v *RuntimeVisitor) VisitGreater(n *parser.NodeGreater) (error) {
    
    return nil
}

func (v *RuntimeVisitor) VisitSmaller(n *parser.NodeSmaller) (error) {
    
    return nil
}

func (v *RuntimeVisitor) VisitOr(n *parser.NodeOr) (error) {
    v.VisitExpr(n.LeftValue())
    boolLeft := boolRegister
    v.VisitExpr(n.RightValue())
    boolRight := boolRegister

    boolRegister = boolLeft || boolRight 
    return nil
} 

func (v *RuntimeVisitor) VisitAnd(n *parser.NodeAnd) (error) {
    v.VisitExpr(n.LeftValue())
    boolLeft := boolRegister
    v.VisitExpr(n.RightValue())
    boolRight := boolRegister

    boolRegister = boolLeft && boolRight 
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
            "author_email",
            "committer",
            "committer_email",
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

    openRepository(path)

    gb.walk, _ = repo.Walk()
    gb.walk.PushHead()
    gb.walk.Sorting(git.SortTime)

    return gb
}



func openRepository(path *string) {
    _repo, err := git.OpenRepository(*path)
    if err != nil {
        panic(err)
    }
    repo = _repo
}

func (g *GitBuilder) setObject(object *git.Commit) {
    g.object = object
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
    if field == "*" {
        return nil
    }

    err := g.isValidTable(tableName)
    if err != nil {
        return err
    }

    table := g.possibleTables[tableName]
    for _, t := range table {
        if t == field {
            return nil
        }
    }

    return throwRuntimeError(fmt.Sprintf("Table '%s' has not field '%s'", tableName, field), 0)
}

// Criar varias funcoes de asserção, a closure usará elas para saber se um certo objeto
// pode ser mostrado ou não.