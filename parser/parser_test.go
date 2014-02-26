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
    New("select")
    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    if ast.child == nil {
        t.Errorf("Program is empty")
    }
}

func TestUsingWildCard(t *testing.T) {
    New("select *")
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