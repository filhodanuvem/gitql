package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	_ "unicode"

	"github.com/cloudson/gitql/lexical"
)

var look_ahead uint8

const Time_YMD = "2006-01-02"
const Time_YMDHIS = "2006-01-02 15:04:05"

type SyntaxError struct {
	expected uint8
	found    uint8
}

func (e *SyntaxError) Error() string {
	var appendix = ""
	if e.found == lexical.T_LITERAL || e.found == lexical.T_ID {
		appendix = fmt.Sprintf("(%s)", lexical.CurrentLexeme)
	}
	return fmt.Sprintf("Expected %s and found %s%s", lexical.TokenName(e.expected), lexical.TokenName(e.found), appendix)
}

func throwSyntaxError(expectedToken uint8, foundToken uint8) error {
	error := new(SyntaxError)
	error.expected = expectedToken
	error.found = foundToken

	return error
}

func New(source string) {
	lexical.New(source)
}

func AST() (*NodeProgram, error) {
	program := new(NodeProgram)
	node, err := gProgram()
	program.Child = node

	return program, err
}

func gProgram() (NodeMain, error) {
	token, tokenError := lexical.Token()
	look_ahead = token

	if tokenError != nil {
		return nil, tokenError
	}

	s, err := gSelect()
	if s == nil {
		return nil, err
	}

	if look_ahead != lexical.T_EOF {
		return nil, throwSyntaxError(lexical.T_EOF, look_ahead)
	}

	return s, nil
}

func gSelect() (*NodeSelect, error) {
	if look_ahead != lexical.T_SELECT {
		return nil, throwSyntaxError(lexical.T_SELECT, look_ahead)
	}
	token, tokenError := lexical.Token()

	look_ahead = token
	if tokenError != nil {
		return nil, tokenError
	}
	s := new(NodeSelect)

	// PARAMETERS
	distincts, fields, err := gTableParams()
	if err != nil {
		return nil, err
	}

	if len(fields) == 1 {
		f0 := fields[0]
		if f0 == "*" {
			s.WildCard = true
		} else if f0 == "#" {
			s.Count = true
		}
	}
	s.Distinct = distincts
	s.Fields = fields

	// TABLES
	tables, err := gTableNames()
	if err != nil {
		return nil, err
	}
	s.Tables = tables

	// WHERE
	where, err := gWhere()
	if err != nil {
		return nil, err
	}
	s.Where = where

	// ORDER BY
	order, err := gOrder()
	if err != nil {
		return nil, err
	}
	s.Order = order

	// LIMIT
	s.Limit, err = gLimit()
	if s.Limit == -1 {
		// @todo search default limit from file config
		s.Limit = 10
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}

func gTableNames() ([]string, error) {
	if look_ahead != lexical.T_FROM {
		return nil, throwSyntaxError(lexical.T_FROM, look_ahead)
	}
	token, error := lexical.Token()
	if error != nil {
		return nil, error
	}
	look_ahead = token
	if look_ahead != lexical.T_ID {
		return nil, throwSyntaxError(lexical.T_ID, look_ahead)
	}

	tables := make([]string, 1)
	tables[0] = lexical.CurrentLexeme

	token2, err := lexical.Token()
	if err != nil && token2 != lexical.T_EOF {
		return nil, err
	}
	look_ahead = token2

	return tables, nil
}

func gTableParams() ([]string, []string, error) {
	if look_ahead == lexical.T_WILD_CARD {
		token, err := lexical.Token()
		if err != nil {
			return nil, nil, err
		}
		look_ahead = token
		return nil, []string{"*"}, nil
	} else if look_ahead == lexical.T_COUNT {
		result, err := gCount()
		return nil, result, err
	}

	isDistinct := false
	if isDistinct = look_ahead == lexical.T_DISTINCT; isDistinct {
		token, err := lexical.Token()
		if err != nil {
			return nil, nil, err
		}
		look_ahead = token
	}
	var fields = []string{}
	var distincts = []string{}
	if look_ahead == lexical.T_ID {
		fields := append(fields, lexical.CurrentLexeme)
		if isDistinct {
			distincts = append(distincts, lexical.CurrentLexeme)
		}
		token, err := lexical.Token()
		if err != nil {
			return nil, nil, err
		}
		look_ahead = token
		fields, errorSyntax := gTableParamsRest(&fields, 1)

		return distincts, fields, errorSyntax
	}
	return nil, nil, throwSyntaxError(lexical.T_ID, look_ahead)
}

// consume count(*)
func gCount() ([]string, error) {
	// by construction, T_COUNT is consumed and stored
	// in the look_ahead
	err := gExactlyASpecificToken(lexical.T_COUNT)
	if err != nil {
		return nil, err
	}
	err = gExactlyASpecificToken(lexical.T_PARENTH_L)
	if err != nil {
		return nil, err
	}
	err = gExactlyASpecificToken(lexical.T_WILD_CARD)
	if err != nil {
		return nil, err
	}
	err = gExactlyASpecificToken(lexical.T_PARENTH_R)
	if err != nil {
		return nil, err
	}
	return []string{"#"}, nil
}

func gExactlyASpecificToken(expected uint8) error {
	if look_ahead != expected {
		return throwSyntaxError(expected, look_ahead)
	}
	token, err := lexical.Token()
	if err != nil {
		return err
	}
	look_ahead = token
	return nil
}

func gTableParamsRest(fields *[]string, count int) ([]string, error) {
	if lexical.T_COMMA == look_ahead {
		var errorToken *lexical.TokenError
		look_ahead, errorToken = lexical.Token()
		if errorToken != nil {
			return *fields, errorToken
		}
		if look_ahead != lexical.T_ID {
			return *fields, throwSyntaxError(lexical.T_ID, look_ahead)
		}

		n := append(*fields, lexical.CurrentLexeme)
		fields = &n
		look_ahead, errorToken = lexical.Token()
		if errorToken != nil {
			return *fields, errorToken
		}
		n, errorSyntax := gTableParamsRest(fields, count+1)
		fields = &n
		if errorSyntax != nil {
			return *fields, errorSyntax
		}
	}

	return *fields, nil
}

func gOrder() (*NodeOrder, error) {
	if look_ahead == lexical.T_ORDER {
		token, err := lexical.Token()
		if err != nil {
			return nil, err
		}
		if token != lexical.T_BY {
			return nil, throwSyntaxError(lexical.T_BY, token)
		}

		order := new(NodeOrder)
		token, err = lexical.Token()
		if err != nil {
			return nil, err
		}
		if token != lexical.T_ID {
			return nil, throwSyntaxError(lexical.T_ID, token)
		}
		order.Field = lexical.CurrentLexeme
		token, err = lexical.Token()
		if err != nil {
			return nil, err
		}
		if token != lexical.T_ASC && token != lexical.T_DESC {
			return nil, throwSyntaxError(lexical.T_ASC, token)
		}
		order.Asc = (token == lexical.T_ASC)

		token, err = lexical.Token()
		if err != nil {
			return nil, err
		}
		look_ahead = token
		return order, nil
	}

	return nil, nil
}

func gLimit() (int, error) {
	if look_ahead != lexical.T_LIMIT {
		return -1, nil
	}
	token, err := lexical.Token()
	if err != nil {
		return 0, err
	}
	look_ahead = token

	number, numberError := strconv.Atoi(lexical.CurrentLexeme)
	if numberError != nil {
		return 0, numberError
	}
	token2, err := lexical.Token()
	if token2 != lexical.T_EOF && err != nil {
		return 0, err
	}
	look_ahead = token2

	return number, nil
}

func gWhere() (NodeExpr, error) {
	if look_ahead != lexical.T_WHERE {
		return nil, nil
	}

	token, tokenError := lexical.Token()
	if tokenError != nil {
		return nil, tokenError
	}
	look_ahead = token
	conds, err := gWhereConds()

	return conds, err
}

func gWhereConds() (NodeExpr, error) {
	where, err := gWC2(false)
	if err != nil {
		return nil, err
	}

	return where, nil
}

// where cond OR
func gWC2(eating bool) (NodeExpr, error) {
	if eating {
		token, err := lexical.Token()
		if token != lexical.T_EOF && err != nil {
			return nil, err
		}
		look_ahead = token
	}
	expr, err := gWC3(false)
	if err != nil {
		return nil, err
	}
	if look_ahead == lexical.T_OR {
		or := new(NodeOr)
		or.SetLeftValue(expr)
		expr2, err := gWC2(true)
		if err != nil {
			return nil, err
		}
		or.SetRightValue(expr2)
		return or, nil
	}
	return expr, nil
}

// where cond AND
func gWC3(eating bool) (NodeExpr, error) {
	if eating {
		token, err := lexical.Token()
		if token != lexical.T_EOF && err != nil {
			return nil, err
		}
		look_ahead = token
	}
	expr, err := gWC4(false)
	if err != nil {
		return nil, err
	}
	if look_ahead == lexical.T_AND {
		and := new(NodeAnd)
		and.SetLeftValue(expr)
		expr2, err := gWC3(true)
		if err != nil {
			return nil, err
		}
		and.SetRightValue(expr2)
		return and, nil
	}
	return expr, nil
}

// where cond 'equal', 'in', 'not in', 'like', 'not like' and 'not equal'
func gWC4(eating bool) (NodeExpr, error) {
	if eating {
		token, err := lexical.Token()
		if token != lexical.T_EOF && err != nil {
			return nil, err
		}
		look_ahead = token
	}
	expr, err := gWC5(false)
	if err != nil {
		return nil, err
	}

	var notBool bool
	if look_ahead == lexical.T_NOT {
		notBool = true
		token, err := lexical.Token()
		if err != nil {
			return nil, err
		}
		look_ahead = token
		if look_ahead != lexical.T_LIKE && look_ahead != lexical.T_IN {
			return nil, throwSyntaxError(lexical.T_NOT, look_ahead)
		}
	}

	switch look_ahead {
	case lexical.T_EQUAL:
		op := new(NodeEqual)
		op.SetLeftValue(expr)
		expr2, err := gWC4(true)
		if err != nil {
			return nil, err
		}
		op.SetRightValue(expr2)

		return op, nil
	case lexical.T_NOT_EQUAL:
		op := new(NodeNotEqual)
		op.SetLeftValue(expr)
		expr2, err := gWC4(true)
		if err != nil {
			return nil, err
		}
		op.SetRightValue(expr2)
		return op, nil
	case lexical.T_LIKE:
		op := new(NodeLike)
		op.SetLeftValue(expr)
		expr2, err := gWC4(true)
		if err != nil {
			return nil, err
		}
		op.SetRightValue(expr2)
		// Compile the regex while parsing, so that
		// we don't need to compile for every row
		rx := strings.Replace(expr2.(*NodeLiteral).Value(), "%", "(.*)", -1)
		op.Pattern, err = regexp.Compile(rx)
		op.Not = notBool
		return op, err
	case lexical.T_IN:
		op := new(NodeIn)
		op.SetLeftValue(expr)
		expr2, err := gWC4(true)
		if err != nil {
			return nil, err
		}
		op.SetRightValue(expr2)
		op.Not = notBool
		return op, nil
	}

	return expr, nil
}

// where cond greater and lesser
func gWC5(eating bool) (NodeExpr, error) {
	if eating {
		token, err := lexical.Token()

		if token != lexical.T_EOF && err != nil {
			return nil, err
		}
		look_ahead = token
	}
	expr, err := rValue()
	if err != nil {
		return nil, err
	}

	switch look_ahead {
	case lexical.T_GREATER, lexical.T_GREATER_OR_EQUAL:
		op := new(NodeGreater)
		op.Equal = (look_ahead == lexical.T_GREATER_OR_EQUAL)
		op.SetLeftValue(expr)
		expr2, err := gWC5(true)
		if err != nil {
			return nil, err
		}
		op.SetRightValue(expr2)

		return op, nil
	case lexical.T_SMALLER, lexical.T_SMALLER_OR_EQUAL:
		op := new(NodeSmaller)
		op.Equal = (look_ahead == lexical.T_SMALLER_OR_EQUAL)
		op.SetLeftValue(expr)
		expr2, err := gWC5(true)
		if err != nil {
			return nil, err
		}
		op.SetRightValue(expr2)

		return op, nil
	}

	return expr, nil
}

func rValue() (NodeExpr, error) {
	if look_ahead == lexical.T_PARENTH_L {
		token, tokenError := lexical.Token()
		if tokenError != nil {
			return nil, tokenError
		}
		look_ahead = token
		conds, err := gWhereConds()
		if err != nil {
			return nil, err
		}
		if look_ahead != lexical.T_PARENTH_R {
			return nil, throwSyntaxError(lexical.T_PARENTH_R, look_ahead)
		}
		token2, tokenError := lexical.Token()
		if token2 != lexical.T_EOF && tokenError != nil {
			return nil, tokenError
		}
		look_ahead = token2

		return conds, nil
	}

	if look_ahead == lexical.T_ID {
		n := new(NodeId)
		n.SetValue(lexical.CurrentLexeme)

		token2, err := lexical.Token()
		if token2 != lexical.T_EOF && err != nil {
			return nil, err
		}
		look_ahead = token2

		return n, nil
	}

	lexeme := lexical.CurrentLexeme
	if look_ahead != lexical.T_LITERAL {
		return nil, errors.New("Only Literals and Date are allowed in `where` clause")
	}

	// @todo inserts IS NULL!

	n := new(NodeLiteral)
	n.SetValue(lexeme)
	token2, err := lexical.Token()
	if err != nil && token2 != lexical.T_EOF {
		return nil, err
	}
	look_ahead = token2

	return n, nil
}
