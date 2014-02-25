package parser

import (
    "testing"
    "reflect"
)

func setUp(source string) {
    
}


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