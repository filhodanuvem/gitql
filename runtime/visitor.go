package runtime

import (
	"reflect"

	"github.com/cloudson/gitql/parser"
)

func (v *RuntimeVisitor) Visit(n *parser.NodeProgram) error {
	return v.VisitSelect(n.Child.(*parser.NodeSelect))
}

func (v *RuntimeVisitor) VisitSelect(n *parser.NodeSelect) error {
	if builder.isProxyTable(n.Tables[0]) {
		proxyTableName := n.Tables[0]
		// refactor tree
		proxy := builder.proxyTables[proxyTableName]
		if !n.WildCard {
			err := testAllFieldsFromTable(n.Fields, proxyTableName)
			if err != nil {
				return err
			}

		} else {
			n.Fields = builder.possibleTables[proxyTableName]
			n.WildCard = false
		}

		n.Tables[0] = proxy.table
		var from, to string
		for from, to = range proxy.fields {
			break
		}

		oldWhere := n.Where
		where := new(parser.NodeAnd)
		condition := new(parser.NodeEqual)
		conditionL := new(parser.NodeId)
		conditionL.SetValue(from)
		conditionR := new(parser.NodeLiteral)
		conditionR.SetValue(to)
		condition.SetLeftValue(conditionL)
		condition.SetRightValue(conditionR)

		where.SetLeftValue(condition)
		where.SetRightValue(oldWhere)

		n.Where = where
	}

	table := n.Tables[0]

	var err error
	err = builder.WithTable(table, table)
	if err != nil {
		return err
	}
	return testAllFieldsFromTable(n.Fields, table)
	// Why not visit expression right now ?
	// Because we will, at first, discover the current object
}

func testAllFieldsFromTable(fields []string, table string) error {
	for _, f := range fields {
		err := builder.UseFieldFromTable(f, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *RuntimeVisitor) VisitExpr(n parser.NodeExpr) error {
	switch reflect.TypeOf(n) {
	case reflect.TypeOf(new(parser.NodeEqual)):
		g := n.(*parser.NodeEqual)
		return v.VisitEqual(g)
	case reflect.TypeOf(new(parser.NodeGreater)):
		g := n.(*parser.NodeGreater)
		return v.VisitGreater(g)
	case reflect.TypeOf(new(parser.NodeSmaller)):
		g := n.(*parser.NodeSmaller)
		return v.VisitSmaller(g)
	case reflect.TypeOf(new(parser.NodeOr)):
		g := n.(*parser.NodeOr)
		return v.VisitOr(g)
	case reflect.TypeOf(new(parser.NodeAnd)):
		g := n.(*parser.NodeAnd)
		return v.VisitAnd(g)
	case reflect.TypeOf(new(parser.NodeIn)):
		g := n.(*parser.NodeIn)
		return v.VisitIn(g)
	case reflect.TypeOf(new(parser.NodeNotEqual)):
		g := n.(*parser.NodeNotEqual)
		return v.VisitNotEqual(g)
	}

	return nil
}

func (v *RuntimeVisitor) VisitEqual(n *parser.NodeEqual) error {
	lvalue := n.LeftValue().(*parser.NodeId).Value()
	rvalue := n.RightValue().(*parser.NodeLiteral).Value()
	boolRegister = n.Assertion(metadata(lvalue), rvalue)
	return nil
}

func (v *RuntimeVisitor) VisitNotEqual(n *parser.NodeNotEqual) error {
	lvalue := n.LeftValue().(*parser.NodeId).Value()
	rvalue := n.RightValue().(*parser.NodeLiteral).Value()
	boolRegister = n.Assertion(metadata(lvalue), rvalue)
	return nil
}

func (v *RuntimeVisitor) VisitGreater(n *parser.NodeGreater) error {
	lvalue := n.LeftValue().(*parser.NodeId).Value()
	lvalue = metadata(lvalue)
	rvalue := n.RightValue().(*parser.NodeLiteral).Value()

	boolRegister = n.Assertion(lvalue, rvalue)

	return nil
}

func (v *RuntimeVisitor) VisitSmaller(n *parser.NodeSmaller) error {
	lvalue := n.LeftValue().(*parser.NodeId).Value()
	lvalue = metadata(lvalue)
	rvalue := n.RightValue().(*parser.NodeLiteral).Value()

	boolRegister = n.Assertion(lvalue, rvalue)

	return nil
}

func (v *RuntimeVisitor) VisitOr(n *parser.NodeOr) error {
	v.VisitExpr(n.LeftValue())
	boolLeft := boolRegister
	v.VisitExpr(n.RightValue())
	boolRight := boolRegister

	boolRegister = boolLeft || boolRight
	return nil
}

func (v *RuntimeVisitor) VisitAnd(n *parser.NodeAnd) error {
	v.VisitExpr(n.LeftValue())
	boolLeft := boolRegister
	v.VisitExpr(n.RightValue())
	boolRight := boolRegister

	boolRegister = boolLeft && boolRight
	return nil
}

func (v *RuntimeVisitor) VisitIn(n *parser.NodeIn) error {
	lvalue := n.LeftValue().(*parser.NodeLiteral).Value()
	rvalue := n.RightValue().(*parser.NodeId).Value()
	boolRegister = n.Assertion(lvalue, metadata(rvalue))

	return nil
}

func (v *RuntimeVisitor) Builder() *GitBuilder {
	return nil
}
