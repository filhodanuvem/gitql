package semantical

import (
	"github.com/cloudson/gitql/parser"
)

type Visitor interface {
	Visit(*parser.NodeProgram) error
	VisitSelect(*parser.NodeSelect) error
	VisitExpr(*parser.NodeExpr) error
	VisitGreater(*parser.NodeGreater) error
	VisitSmaller(*parser.NodeSmaller) error
	VisitIn(*parser.NodeSmaller) error
	VisitEqual(*parser.NodeSmaller) error
	VisitNotEqual(*parser.NodeSmaller) error
}

type SemanticalVisitor struct {
	Visitor
}
