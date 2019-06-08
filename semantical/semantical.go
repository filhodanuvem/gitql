package semantical

import (
	"fmt"

	"github.com/navigaid/gitql/lexical"
	"github.com/navigaid/gitql/parser"
)

func Analysis(ast *parser.NodeProgram) error {
	semantic := new(SemanticalVisitor)

	return semantic.Visit(ast)
}

type SemanticalError struct {
	err   string
	errNo uint8
}

func throwSemanticalError(err string) error {
	end := new(SemanticalError)
	end.err = err

	return end
}

func (e *SemanticalError) Error() string {
	return e.err
}

func (v *SemanticalVisitor) Visit(n *parser.NodeProgram) error {
	return v.VisitSelect(n.Child.(*parser.NodeSelect))
}

func (v *SemanticalVisitor) VisitSelect(n *parser.NodeSelect) error {

	fields := n.Fields
	fieldsCount := make(map[string]bool)
	for _, field := range fields {
		if fieldsCount[field] {
			return throwSemanticalError(fmt.Sprintf("Field '%s' found many times", field))
		}

		fieldsCount[field] = true
	}

	err := v.VisitExpr(n.Where)
	if err != nil {
		return err
	}

	if 0 == n.Limit {
		return throwSemanticalError("Limit should be greater than zero")
	}

	return nil
}

func (v *SemanticalVisitor) VisitExpr(n parser.NodeExpr) error {
	if n == nil {
		return nil
	}

	switch n.Operator() {
	case lexical.T_GREATER:
		g := n.(*parser.NodeGreater)
		return v.VisitGreater(g)
	case lexical.T_SMALLER:
		g := n.(*parser.NodeSmaller)
		return v.VisitSmaller(g)
	case lexical.T_IN:
		g := n.(*parser.NodeIn)
		return v.VisitIn(g)
	}

	return nil
}

func (v *SemanticalVisitor) VisitGreater(n *parser.NodeGreater) error {
	rVal := n.RightValue()
	if !shouldBeNumericOrDate(rVal) {
		return throwSemanticalError("RValue in Greater should be numeric or a date")
	}

	return nil
}

func (v *SemanticalVisitor) VisitSmaller(n *parser.NodeSmaller) error {
	rVal := n.RightValue()
	if !shouldBeNumericOrDate(rVal) {
		return throwSemanticalError("RValue in Smaller should be numeric or a date")
	}

	return nil
}

func (v *SemanticalVisitor) VisitIn(n *parser.NodeIn) error {
	lval := n.LeftValue()
	if lval.Operator() != lexical.T_LITERAL {
		return throwSemanticalError("LValue at In operator shoud be a literal")
	}

	rval := n.RightValue()
	if rval.Operator() != lexical.T_ID {
		return throwSemanticalError("RValue at In operator shoud be a Identifier")
	}

	return nil
}

func shouldBeNumericOrDate(val parser.NodeExpr) bool {
	// if reflect.TypeOf(val) == reflect.TypeOf(new(parser.NodeNumber)) {
	if val.Operator() == lexical.T_NUMERIC {
		return true
	}

	if val.Operator() == lexical.T_LITERAL {
		date := parser.ExtractDate(val.(*parser.NodeLiteral).Value())
		if date != nil {
			return true
		}
	}

	return false
}
