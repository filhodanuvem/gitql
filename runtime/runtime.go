package runtime

import (
    "fmt"
    "strings"
    "github.com/cloudson/gitql/parser"
    "github.com/libgit2/git2go"
)

const (
    WALK_COMMITS = 1
    WALK_TREES = 2
    WALK_REFERENCES = 3
)

const (
    REFERENCE_TYPE_BRANCH = "branch"
    REFERENCE_TYPE_REMOTE = "remote"
    REFERENCE_TYPE_TAG = "tag"
)

var repo *git.Repository
var builder *GitBuilder
var boolRegister bool

type GitBuilder struct {
    tables map[string]string 
    possibleTables map[string][]string
    repo *git.Repository
    currentWalkType uint8
    currentCommit *git.Commit
    currentReference *git.Reference
    walk *git.RevWalk

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
    switch findWalkType(n) {
        case WALK_COMMITS: 
            walkCommits(n, visitor)
            break
        case WALK_TREES:
            walkTrees(n, visitor)
            break
        case WALK_REFERENCES:
            walkReferences(n, visitor)
    }
}

func findWalkType(n *parser.NodeProgram) uint8 {
    s := n.Child.(*parser.NodeSelect)

    switch s.Tables[0] {
        case "commits" :
            builder.currentWalkType = WALK_COMMITS
        case "trees" :
            builder.currentWalkType = WALK_TREES
        case "refs":
            builder.currentWalkType = WALK_REFERENCES
    }
    
    return builder.currentWalkType
}

func walkCommits(n *parser.NodeProgram, visitor *RuntimeVisitor) {
    builder.walk, _ = repo.Walk()
    builder.walk.PushHead()
    builder.walk.Sorting(git.SortTime)

    s := n.Child.(*parser.NodeSelect)
    where := s.Where
    
    counter := 1
    fmt.Println()
    fn := func (object *git.Commit) bool {
        builder.setCommit(object)
        boolRegister = true
        visitor.VisitExpr(where)
        if boolRegister {
            fields := s.Fields
            if s.WildCard {
                fields = builder.possibleTables[s.Tables[0]]
            }            
            for _, f := range fields {
                fmt.Printf("%s | ", metadataCommit(f, object))    
            }
            fmt.Println()

            
            counter = counter + 1
        }
        if counter > s.Limit {
            return false
        }
        return true
    }

    err := builder.walk.Iterate(fn)
    if err != nil {
        fmt.Printf(err.Error())
    }
}

func walkTrees(n *parser.NodeProgram, visitor *RuntimeVisitor) {
    // not yet!
}

func walkReferences(n *parser.NodeProgram, visitor *RuntimeVisitor) {
    s := n.Child.(*parser.NodeSelect)
    where := s.Where

    // @TODO make PR with Repository.WalkReference()
    iterator, err := builder.repo.NewReferenceIterator()  
    if err != nil {
        panic(err)
    }
    counter := 1
    for object, inTheEnd := iterator.Next(); inTheEnd == nil; object, inTheEnd = iterator.Next() {
        
        builder.setReference(object)
        boolRegister = true
        visitor.VisitExpr(where)
        if boolRegister {
            fields := s.Fields
            if s.WildCard {
                fields = builder.possibleTables[s.Tables[0]]
            }            
            for _, f := range fields {
                fmt.Printf("%s | ", metadataReference(f, object))    
            }
            fmt.Println()

            counter = counter + 1
            if counter > s.Limit {
                break
            }
        }
    }

}

func metadata(identifier string) string {
    switch builder.currentWalkType {
        case WALK_COMMITS:
            return metadataCommit(identifier, builder.currentCommit)
        case WALK_REFERENCES:
            return metadataReference(identifier, builder.currentReference)
    }

    panic("GOD!")
}

func metadataTree(identifier string, object *git.TreeEntry) string {
    return "" // not yet implemented!
}

func metadataReference(identifier string, object *git.Reference) string {
    switch identifier {
        case "name":
            return object.Shorthand()
        case "full_name" : 
            return object.Name()
        case "hash" :
            target := object.Target()
            if target == nil {
                return "NULL"
            }
            return target.String()
        case "type":
            if object.IsBranch() {
                return REFERENCE_TYPE_BRANCH
            }

            if object.IsRemote() {
                return REFERENCE_TYPE_REMOTE
            }

            if object.IsTag() {
                return REFERENCE_TYPE_TAG
            }
    }

    panic(fmt.Sprintf("Trying select field %s ", identifier))
}

func metadataCommit(identifier string, object *git.Commit) string {
    key := "" 
    for key, _ = range builder.tables {
        break
    }
    table := key
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
            return object.Committer().When.Format(parser.Time_YMDHIS)
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
        "trees": {
            "hash",
            "name",
            "id",
            "type",
            "filemode",
        },
        "refs": {
            "name",
            "full_name",
            "type",
            "hash",
        },
    }
    gb.possibleTables = possibleTables

    openRepository(path)

    gb.repo = repo

    return gb
}



func openRepository(path *string) {
    _repo, err := git.OpenRepository(*path)
    if err != nil {
        panic(err)
    }
    repo = _repo
}

func (g *GitBuilder) setCommit(object *git.Commit) {
    g.currentCommit = object
}

func (g *GitBuilder) setReference(object *git.Reference) {
    g.currentReference = object
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