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


    // if error != nil {
    //     t.Errorf(error.Error())
    // }

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