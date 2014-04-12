package runtime

import (
    "reflect"
    "github.com/cloudson/gitql/parser"
)

func (v *RuntimeVisitor) Visit(n *parser.NodeProgram) (error) {
    return v.VisitSelect(n.Child.(*parser.NodeSelect))
} 

func (v *RuntimeVisitor) VisitSelect(n *parser.NodeSelect) (error) {
    table := n.Tables[0]
    fields := n.Fields 
    var err error
    for _, f := range fields {
        err = builder.UseFieldFromTable(f, table)
        builder.tables[table] = table
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
        case reflect.TypeOf(new(parser.NodeIn)):
            g:= n.(*parser.NodeIn)
            return v.VisitIn(g)

    } 

    return nil
}

func (v *RuntimeVisitor) VisitEqual(n *parser.NodeEqual) (error) {
    lvalue := n.LeftValue().(*parser.NodeId).Value()
    rvalue := n.RightValue().(*parser.NodeLiteral).Value()
    boolRegister = n.Assertion(metadata(lvalue), rvalue)
    
    return nil
}

func (v *RuntimeVisitor) VisitGreater(n *parser.NodeGreater) (error) {
    lvalue := n.LeftValue().(*parser.NodeId).Value()
    lvalue = metadata(lvalue)
    rvalue := n.RightValue().(*parser.NodeLiteral).Value()
    
    boolRegister = n.Assertion(lvalue, rvalue)

    return nil
}

func (v *RuntimeVisitor) VisitSmaller(n *parser.NodeSmaller) (error) {
    lvalue := n.LeftValue().(*parser.NodeId).Value()
    lvalue = metadata(lvalue)
    rvalue := n.RightValue().(*parser.NodeLiteral).Value()

    boolRegister = n.Assertion(lvalue, rvalue)
    
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

func (v *RuntimeVisitor) VisitIn(n *parser.NodeIn) (error) {
    lvalue := n.LeftValue().(*parser.NodeLiteral).Value()
    rvalue := n.RightValue().(*parser.NodeId).Value()
    boolRegister = n.Assertion(lvalue, metadata(rvalue))

    return nil
}  

func (v *RuntimeVisitor) Builder() (*GitBuilder){
    return nil
}
