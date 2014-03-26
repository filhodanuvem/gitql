package semantical

import (
    "testing"
    "github.com/cloudson/gitql/parser"
)

func TestInvalidZeroLimit(t *testing.T) {
    parser.New("select * from commits limit 0")
    ast, parserErr := parser.AST()
    if parserErr != nil {
        t.Fatalf(parserErr.Error())
    }

    err := analysis(ast)
    if err == nil {
        t.Fatalf("Should not accept limit zero")
    }
}

func TestValidNullLimit(t *testing.T) {
    parser.New("select * from commits")
    ast, parserErr := parser.AST()
    if parserErr != nil {
        t.Fatalf(parserErr.Error())
    }

    err := analysis(ast)
    if err != nil {
        t.Fatalf(err.Error())
    }
}

func TestChooseRepetitiveFields(t *testing.T) {
    parser.New("select name, created_at, name from commits")
    ast, parserErr := parser.AST()
    if parserErr != nil {
        t.Fatalf(parserErr.Error())
    }

    err := analysis(ast)
    if err == nil {
        t.Fatalf("Shoud avoid repetitive fields")
    }
}

func TestConstantLValue(t *testing.T) {
    parser.New("select name from commits where 'name' = 'name' ")
    ast, parserErr := parser.AST()
    if parserErr != nil {
        t.Fatalf(parserErr.Error())
    }

    err := analysis(ast)
    if err != nil {
        t.Fatalf(err.Error())
    }
}

func TestGreaterWithNoNumeric(t *testing.T) {
    parser.New("select name from commits where date > 'name'")
    ast, parserErr := parser.AST()
    if parserErr != nil {
        t.Fatalf(parserErr.Error())
    }

    err := analysis(ast)
    if err == nil {
        t.Fatalf("Shoud avoid greater with no numeric")
    }
}

func TestSmallerWithNoNumeric(t *testing.T) {
    parser.New("select name from commits where date > 'name'")
    ast, parserErr := parser.AST()
    if parserErr != nil {
        t.Fatalf(parserErr.Error())
    }

    err := analysis(ast)
    if err == nil {
        t.Fatalf("Shoud avoid smaller with no numeric")
    }
}