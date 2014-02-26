package parser

import (
    "fmt"
    "github.com/cloudson/gitql"
    "github.com/cloudson/gitql/lexical"
)

type SyntaxError struct {
    expected uint8
    found uint8
}

func (e *SyntaxError) Error() string {
    var appendix = ""; 
    if e.found == lexical.T_LITERAL {
        appendix = fmt.Sprintf("(%s)", lexical.CurrentLexeme)
    }
    return fmt.Sprintf("Expected %d and found %d%s", e.expected, e.found, appendix)
}

func (s *NodeSelect) Run() {
    return 
}

func throwSyntaxError(expectedToken uint8, foundToken uint8) (gitql.CompileError){
    error := new(SyntaxError)
    error.expected = expectedToken
    error.found = foundToken

    return error
}

func New(source string) {
    lexical.New(source)
}

func AST() (p *NodeProgram, error gitql.CompileError) {
    program := new(NodeProgram)
    program.child, error = g_program()

    return program, error
}

func g_program() (NodeMain, gitql.CompileError) {
    token, _ := lexical.Token()
    if token != lexical.T_SELECT {
        return nil, throwSyntaxError(lexical.T_SELECT, token)
    }

    s := new(NodeSelect)
    s.WildCard = true

    return s, nil
}

func g_table_params() ([]string, gitql.CompileError){
    token, error := lexical.Token()
    if error != nil {
        return nil, error
    }

    if token != lexical.T_LITERAL && token != lexical.T_WILD_CARD {
        return nil, throwSyntaxError(lexical.T_LITERAL, token)
    }

    if token == lexical.T_WILD_CARD {
        return []string{"*"}, nil
    }

    return g_table_params_rest()
}

func g_table_params_rest() ([]string, gitql.CompileError){
    return nil, nil
}