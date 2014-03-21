package parser

import "github.com/cloudson/gitql/lexical"
import "strconv"

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
    Operator() uint8
    LeftValue() NodeExpr
    RightValue() NodeExpr
    SetLeftValue(NodeExpr) 
    SetRightValue(NodeExpr) 
}

type NodeEqual struct {
    leftValue NodeExpr
    rightValue NodeExpr
}

type NodeAnd struct { 
    leftValue NodeExpr
    rightValue NodeExpr
}

type NodeNumber struct {
    value float64
    leftValue NodeExpr
    rightValue NodeExpr
}

type NodeLiteral struct {
    leftValue NodeExpr
    rightValue NodeExpr
    value string
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

func (n *NodeEqual) Operator() uint8{
    return lexical.T_EQUAL
}

func (n *NodeEqual) SetLeftValue(e NodeExpr) {
    n.leftValue = e
}

func (n *NodeEqual) SetRightValue(e NodeExpr) {
    n.rightValue = e
}

func (n *NodeEqual) RightValue() NodeExpr{
    return n.rightValue
}

func (n *NodeEqual) LeftValue() NodeExpr{
    return n.leftValue
}

func (n *NodeNumber) Operator() uint8{
    return 0
}

func (n *NodeLiteral) Operator() uint8{
    return lexical.T_LITERAL
}

func (n *NodeLiteral) SetValue(value string) {
    n.value = value
}

func (n *NodeLiteral) SetLeftValue(e NodeExpr) {
    n.leftValue = e
}

func (n *NodeLiteral) SetRightValue(e NodeExpr) {
    n.rightValue = e
}

func (n *NodeLiteral) RightValue() NodeExpr{
    return n.rightValue
}

func (n *NodeLiteral) LeftValue() NodeExpr{
    return n.leftValue
}

func (n *NodeNumber) SetValue(value string) {
    n.value, _ = strconv.ParseFloat(value, 64)
}

func (n *NodeNumber) SetLeftValue(e NodeExpr) {
    n.leftValue = e
}

func (n *NodeNumber) SetRightValue(e NodeExpr) {
    n.rightValue = e
}

func (n *NodeNumber) RightValue() NodeExpr{
    return n.rightValue
}

func (n *NodeNumber) LeftValue() NodeExpr{
    return n.leftValue
}
