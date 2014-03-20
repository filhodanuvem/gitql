package parser

import (
    "testing"
    "reflect"
)


func TestEmptySource(t *testing.T) {
    New("")
    ast, _ := AST()

    if (reflect.TypeOf(ast) != reflect.TypeOf(new(NodeProgram))) {
        t.Errorf("AST should be a NodeProgram, found %s", reflect.TypeOf(ast).String())
    }

    if (ast.child != nil) {
        t.Errorf("Program is not empty")
    }
}

func TestInvalidFirstNode(t *testing.T) {
    New("cloudson")
    _, error := AST()

    if error == nil {
        t.Errorf("Expected a syntax error")
    }
}

func TestValidFirstNode(t *testing.T) {
    New("select * from users")
    ast, _ := AST()

    if ast.child == nil {
        t.Errorf("Program is empty")
    }
}

func TestUsingWildCard(t *testing.T) {
    New("select * from users")
    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    if ast.child == nil {
        t.Errorf("Program is empty")
    }

    selectNode := ast.child.(*NodeSelect)
    if !selectNode.WildCard {
        t.Errorf("Expected wildcard setted")
    }
}

func TestUsingOneFieldName(t *testing.T) {
    New("select name from files")

    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.child.(*NodeSelect)
    
    if len(selectNode.fields) != 1 {
        t.Errorf("Expected exactly one field and found %d", len(selectNode.fields))
    }

    if selectNode.fields[0] != "name" {
        t.Errorf("Expected param 'name' and found '%s'", selectNode.fields[0])
    }
}

func TestUsingFieldNames(t *testing.T) {
    New("select name, created_at from files")

    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.child.(*NodeSelect)
    if len(selectNode.fields) != 2 {
        t.Errorf("Expected exactly two fields and found %d", len(selectNode.fields))
    }
}

func TestWithOneTable(t *testing.T) {
    New("select name, created_at from files")

    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.child.(*NodeSelect)
    if len(selectNode.fields) != 2 {
        t.Errorf("Expected exactly two fields and found %d", len(selectNode.fields))
    }

    if selectNode.tables[0] != "files" {
        t.Errorf("Expected table 'files', found %s", selectNode.tables[0])
    }
}

func TestErrorWithUnexpectedComma(t *testing.T) {
    New("select name, from files")

    _, error := AST()

    if error == nil {
        t.Errorf("Expected error 'Unexpected T_COMMA'")
    }
}

func TestErrorWithInvalidRootNode(t *testing.T) {
    New("name from files")

    _, error := AST()
    if error == nil {
        t.Errorf("Expected error 'EXPECTED T_SELECT'")
    }

}

func testErrorSqlWithoutTable(t *testing.T) {
    New("select name from ")

    _, error := AST()
    if error == nil {
        t.Errorf("Expected error 'EXPECTED table'")   
    }
}

func testWithLimit(t *testing.T) {
    New("select * from files where limit 5")

    ast, error := AST()
    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.child.(*NodeSelect)
    if selectNode.limit != 5 {
        t.Errorf("Limit should be 5!!!")   
    }
}

func testWithEmptyLimit(t *testing.T) {
    New("select * from files where limit")

    _, error := AST()
    if error == nil {
        t.Errorf("Shoud throw error because limit has not value")
    }
}

func testWithNonNumericLimit(t *testing.T) {
    New("select * from commits where limit cloud")

    _, error := AST() 
    if error == nil {
        t.Errorf("Shoud throw error because limit is not a number")
    }
}

func testWithWhereSimpleComparting(t *testing.T) {
    New("select * from commits where hash = e69de29bb2d1d6434b8b29ae775ad8c2e48c5391")

    ast, err := AST() 
    if err != nil {
        t.Errorf(err.Error())
    }

    selectNode := ast.child.(*NodeSelect)
    if selectNode.where == nil{
        t.Errorf("should has where node")
    }
}