package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloudson/gitql/lexical"
)

type NodeMain interface {
	Run()
}

type NodeEmpty struct {
}

type NodeProgram struct {
	Child NodeMain
	Path  *string
}

type NodeSelect struct {
	WildCard bool
	Count    bool
	Distinct []string
	Fields   []string
	Tables   []string
	Where    NodeExpr
	Order    *NodeOrder
	Limit    int
}

type NodeExpr interface {
	Assertion(lvalue, rvalue string) bool
	Operator() uint8
	LeftValue() NodeExpr
	RightValue() NodeExpr
	SetLeftValue(NodeExpr)
	SetRightValue(NodeExpr)
}

type NodeBinOp interface {
	LeftValue() NodeExpr
	RightValue() NodeExpr
	SetLeftValue(NodeExpr)
	SetRightValue(NodeExpr)
}

type NodeConst interface {
	SetValue(string)
}

type NodeAdapterBinToConst struct {
	adapted NodeBinOp
}

type NodeIn struct {
	leftValue  NodeExpr
	rightValue NodeExpr
	Not        bool
}

type NodeEqual struct {
	leftValue  NodeExpr
	rightValue NodeExpr
}

type NodeNotEqual struct {
	leftValue  NodeExpr
	rightValue NodeExpr
}

type NodeLike struct {
	leftValue  NodeExpr
	rightValue NodeExpr
	Pattern    *regexp.Regexp
	Not        bool
}

type NodeGreater struct {
	leftValue  NodeExpr
	rightValue NodeExpr
	Equal      bool
}

type NodeSmaller struct {
	leftValue  NodeExpr
	rightValue NodeExpr
	Equal      bool
}

type NodeOr struct {
	leftValue  NodeExpr
	rightValue NodeExpr
}

type NodeAnd struct {
	leftValue  NodeExpr
	rightValue NodeExpr
}

type NodeNumber struct {
	value      float64
	leftValue  NodeExpr
	rightValue NodeExpr
}

type NodeLiteral struct {
	leftValue  NodeExpr
	rightValue NodeExpr
	value      string
}

type NodeId struct {
	leftValue  NodeExpr
	rightValue NodeExpr
	value      string
}

type NodeOrder struct {
	Field string
	Asc   bool
}

func (s *NodeSelect) Run() {
	return
}

func (e *NodeEmpty) Run() {
	return
}

func (n *NodeIn) Assertion(lvalue string, rvalue string) bool {
	if n.Not {
		return !strings.Contains(rvalue, lvalue)
	}
	return strings.Contains(rvalue, lvalue)
}

func (n *NodeIn) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeIn) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeIn) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeIn) LeftValue() NodeExpr {
	return n.leftValue
}

func (n *NodeIn) Operator() uint8 {
	return lexical.T_IN
}

// EQUAL
func (n *NodeEqual) Assertion(lvalue string, rvalue string) bool {
	if len(lvalue) == 40 {
		return lvalue[0:len(rvalue)] == rvalue
	}
	return lvalue == rvalue
}

func (n *NodeEqual) Operator() uint8 {
	return lexical.T_EQUAL
}

func (n *NodeEqual) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeEqual) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeEqual) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeEqual) LeftValue() NodeExpr {
	return n.leftValue
}

// NOT EQUAL
func (n *NodeNotEqual) Assertion(lvalue string, rvalue string) bool {
	if len(lvalue) == 40 {
		return lvalue[0:len(rvalue)] != rvalue
	}
	return lvalue != rvalue
}

func (n *NodeNotEqual) Operator() uint8 {
	return lexical.T_NOT_EQUAL
}

func (n *NodeNotEqual) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeNotEqual) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeNotEqual) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeNotEqual) LeftValue() NodeExpr {
	return n.leftValue
}

// LIKE
func (n *NodeLike) Assertion(lvalue string, rvalue string) bool {
	if n.Not {
		return !n.Pattern.MatchString(lvalue)
	}
	return n.Pattern.MatchString(lvalue)
}

func (n *NodeLike) Operator() uint8 {
	return lexical.T_LIKE
}

func (n *NodeLike) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeLike) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeLike) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeLike) LeftValue() NodeExpr {
	return n.leftValue
}

// GREATER
func (n *NodeGreater) Assertion(lvalue string, rvalue string) bool {
	time := ExtractDate(rvalue)
	if time != nil {
		timeFound := ExtractDate(lvalue)
		if timeFound != nil {
			return timeFound.After(*time) || (n.Equal && timeFound.Equal(*time))
		}
	}
	return lvalue > rvalue
}

func (n *NodeGreater) Operator() uint8 {
	return lexical.T_GREATER
}

func (n *NodeGreater) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeGreater) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeGreater) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeGreater) LeftValue() NodeExpr {
	return n.leftValue
}

// SMALLER
func (n *NodeSmaller) Assertion(lvalue string, rvalue string) bool {
	time := ExtractDate(rvalue)
	if time != nil {
		timeFound := ExtractDate(lvalue)
		if timeFound != nil {
			return timeFound.Before(*time) || (n.Equal && timeFound.Equal(*time))
		}
	}
	return lvalue < rvalue
}

func (n *NodeSmaller) Operator() uint8 {
	return lexical.T_SMALLER
}

func (n *NodeSmaller) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeSmaller) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeSmaller) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeSmaller) LeftValue() NodeExpr {
	return n.leftValue
}

// OR
func (n *NodeOr) Assertion(lvalue string, rvalue string) bool {
	return false

}

func (n *NodeOr) Operator() uint8 {
	return lexical.T_OR
}

func (n *NodeOr) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeOr) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeOr) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeOr) LeftValue() NodeExpr {
	return n.leftValue
}

// AND
func (n *NodeAnd) Assertion(lvalue string, rvalue string) bool {
	return lvalue == rvalue
}

func (n *NodeAnd) Operator() uint8 {
	return lexical.T_AND
}

func (n *NodeAnd) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeAnd) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeAnd) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeAnd) LeftValue() NodeExpr {
	return n.leftValue
}

// LITERAL
func (n *NodeLiteral) Assertion(lvalue string, rvalue string) bool {
	return lvalue == rvalue
}

func (n *NodeLiteral) Operator() uint8 {
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

func (n *NodeLiteral) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeLiteral) LeftValue() NodeExpr {
	return n.leftValue
}

func (n *NodeLiteral) Value() string {
	return n.value
}

// NUMBER
func (n *NodeNumber) Assertion(lvalue string, rvalue string) bool {
	return lvalue == rvalue
}

func (n *NodeNumber) Operator() uint8 {
	return lexical.T_NUMERIC
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

func (n *NodeNumber) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeNumber) LeftValue() NodeExpr {
	return n.leftValue
}

func (n *NodeNumber) Value() float64 {
	return n.value
}

// ID
func (n *NodeId) Assertion(lvalue string, rvalue string) bool {
	return lvalue == rvalue
}

func (n *NodeId) Operator() uint8 {
	return lexical.T_ID
}

func (n *NodeId) SetValue(value string) {
	n.value = value
}

func (n *NodeId) SetLeftValue(e NodeExpr) {
	n.leftValue = e
}

func (n *NodeId) SetRightValue(e NodeExpr) {
	n.rightValue = e
}

func (n *NodeId) RightValue() NodeExpr {
	return n.rightValue
}

func (n *NodeId) LeftValue() NodeExpr {
	return n.leftValue
}

func (n *NodeId) Value() string {
	return n.value
}

func (n *NodeAdapterBinToConst) setAdapted(a NodeBinOp) {
	n.adapted = a
}

func ExtractDate(date string) *time.Time {
	t, err := time.Parse(Time_YMD, date)
	if err == nil {
		return &t
	}

	t, err = time.Parse(Time_YMDHIS, date)
	if err == nil {
		return &t
	}

	// does not matter if the string is not a date
	// gitql will use it like a simple text
	return nil
}
