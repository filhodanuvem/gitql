package semantical 

import (
    "fmt"
    "github.com/cloudson/gitql/parser"
)

type Visitor interface {
    Visit(*parser.NodeProgram) error 
    VisitSelect(*parser.NodeSelect) error 
}

type SemanticalVisitor struct {

}

type SemanticalError struct {
    err string
}

func throwSemanticalError(err string) (error) {
    end := new(SemanticalError)
    end.err = err 

    return end
}

func (e *SemanticalError) Error() string {
    return e.err
}

func (v *SemanticalVisitor) Visit(n *parser.NodeProgram) (error) {
    return v.VisitSelect(n.Child.(*parser.NodeSelect))
} 

func (v *SemanticalVisitor) VisitSelect(n *parser.NodeSelect) (error) {
    

    fields := n.Fields 
    fieldsCount := make(map[string]bool)
    for _, field := range fields {
        if fieldsCount[field] {
           return throwSemanticalError(fmt.Sprintf("Field '%s' found may times", field))
        }

        fieldsCount[field] = true
    }

    if 0 == n.Limit {
        return throwSemanticalError("Limit should be greater then zero")
    }

    return nil 
} 