package parser

import "github.com/cloudson/gitql/lexical"

type NodeMain interface {
    Run()
}

type NodeEmpty struct {

}

type NodeProgram struct {
    child NodeMain
}

type NodeSelect struct {
    WildCard bool
    fields []string
    tables []string
    where NodeExpr
    limit int 
}

type NodeExpr interface {
    Operator()
}

type NodeAnd struct { 
    LeftValue NodeExpr
    RightValue NodeExpr
}

type NodeNumber struct {
    Value int
}

func (s *NodeSelect) Run() {
    return 
}

func (e *NodeEmpty) Run() {
    return 
}

func (n *NodeAnd) Operator() uint8{
    return lexical.T_AND
}

func (n *NodeNumber) Operator() uint8{
    return 0
}