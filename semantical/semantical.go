package semantical

import (
    "fmt"
    "reflect"
    "github.com/cloudson/gitql/parser"
)

func Analysis(ast *parser.NodeProgram) error {
    semantic := new(SemanticalVisitor)

    return semantic.Visit(ast)
}

type SemanticalError struct {
    err string
    errNo uint8
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


    err := v.VisitExpr(n.Where)
    if err != nil {
        return err
    }

    if 0 == n.Limit {
        return throwSemanticalError("Limit should be greater then zero")
    }

    return nil 
} 

func (v *SemanticalVisitor) VisitExpr(n parser.NodeExpr) (error) {
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

func (v *SemanticalVisitor) VisitGreater(n *parser.NodeGreater) (error) {
    rVal := n.RightValue()
    if !shouldBeNumericOrDate(rVal) {
        return throwSemanticalError("RValue in Greater should be numeric or a date")
    } 

    return nil
}

func (v *SemanticalVisitor) VisitSmaller(n *parser.NodeSmaller) (error) {
    rVal := n.RightValue()
    if !shouldBeNumericOrDate(rVal) {
        return throwSemanticalError("RValue in Smaller should be numeric or a date")
    } 

    return nil
}

func shouldBeNumericOrDate(val parser.NodeExpr) bool {
    if reflect.TypeOf(val) == reflect.TypeOf(new(parser.NodeNumber)) {
        return true
    }

    if reflect.TypeOf(val) == reflect.TypeOf(new(parser.NodeLiteral)) {
        date := parser.ExtractDate(val.(*parser.NodeLiteral).Value())        
        if date != nil {
            return true
        }
    }
    
    return false
}
