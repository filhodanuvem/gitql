package parser

import (
    "github.com/cloudson/gitql/lexical"
)

type SyntaxError struct {
    expected uint8
    found uint8
}

func (e *SyntaxError) Error() string {
    return ""
}

func throwSyntaxError(expectedToken uint8, foundToken uint8) (*SyntaxError){
    error := new(SyntaxError)
    error.expected = expectedToken
    error.found = foundToken

    return error
}


func New(source string) {
    lexical.New(source)
}

func AST() (*NodeProgram, *SyntaxError){
    program := new(NodeProgram)

    return program, nil
}