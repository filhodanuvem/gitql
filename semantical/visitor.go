package semantical 

import (
    "github.com/cloudson/gitql/parser"
)

const Time_YMD = "2006-01-02"
const Time_YMDHIS = "2006-01-02 15:04:05"

type Visitor interface {
    Visit(*parser.NodeProgram) error 
    VisitSelect(*parser.NodeSelect) error 
    VisitExpr(*parser.NodeExpr) error
    VisitGreater(*parser.NodeGreater) error
    VisitSmaller(*parser.NodeSmaller) error
}

type SemanticalVisitor struct {

}

