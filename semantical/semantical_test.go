package semantical

import (
	"testing"

	"github.com/navigaid/gitql/parser"
)

func TestInvalidZeroLimit(t *testing.T) {
	parser.New("select * from commits limit 0")
	ast, parserErr := parser.AST()
	if parserErr != nil {
		t.Fatalf(parserErr.Error())
	}

	err := Analysis(ast)
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

	err := Analysis(ast)
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

	err := Analysis(ast)
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

	err := Analysis(ast)
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

	err := Analysis(ast)
	if err == nil {
		t.Fatalf("Shoud avoid greater with no numeric")
	}
}

func TestSmallerWithInvalidConstant(t *testing.T) {
	parser.New("select name from commits where date <= 'name'")
	ast, parserErr := parser.AST()
	if parserErr != nil {
		t.Fatalf(parserErr.Error())
	}

	err := Analysis(ast)
	if err == nil {
		t.Fatalf("Shoud avoid smaller with no numeric")
	}
}

func TestSmallerWithDate(t *testing.T) {
	parser.New("select name from commits where date > '2013-03-14 00:00:00'")
	ast, parserErr := parser.AST()
	if parserErr != nil {
		t.Fatalf(parserErr.Error())
	}

	err := Analysis(ast)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSmallerWithDateWithoutTime(t *testing.T) {
	parser.New("select count(*) from commits where date > '2013-03-14'")
	ast, parserErr := parser.AST()
	if parserErr != nil {
		t.Fatalf(parserErr.Error())
	}

	err := Analysis(ast)
	if err != nil {
		t.Fatalf(err.Error())
	}

}

// You should not test stupid things like "c" in "cloudson" or 1 = 1 ¬¬
func TestInUsingNotLiteralLeft(t *testing.T) {
	parser.New("select * from commits where 'c' in 'cloudson'")
	ast, parserErr := parser.AST()
	if parserErr != nil {
		t.Fatalf(parserErr.Error())
	}

	err := Analysis(ast)
	if err == nil {
		t.Fatalf("Should trow error with invalid in ")
	}
}

func TestInUsingNotIdRight(t *testing.T) {
	parser.New("select * from commits where 'c' in 'cc' ")
	ast, parserErr := parser.AST()
	if parserErr != nil {
		t.Fatalf(parserErr.Error())
	}

	err := Analysis(ast)
	if err == nil {
		t.Fatalf("Should trow error with invalid in ")
	}
}
