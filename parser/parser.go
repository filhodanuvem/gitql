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
    return fmt.Sprintf("Expected %s and found %s%s", lexical.TokenName(e.expected), lexical.TokenName(e.found), appendix)
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
    params, error := g_table_params()
    if error != nil {
        return nil, error
    }

    if len(params) == 1 && params[0] == "*" {
        s.WildCard = true    
    }
    s.params = params

    return s, nil
}

func g_table_params() ([]string, gitql.CompileError){
    token, errorToken := lexical.Token()
    if errorToken != nil {
        return nil, errorToken
    }

    if token == lexical.T_WILD_CARD {
        return []string{"*"}, nil
    }

    var fields = []string{}
    if token == lexical.T_LITERAL {
        fields := append(fields, lexical.CurrentLexeme)

        fields, errorSyntax := g_table_params_rest(&fields, 1)

        return fields, errorSyntax
    }
    return nil, throwSyntaxError(lexical.T_LITERAL, token)
    
}

func g_table_params_rest(fields *[]string, count int) ([]string, gitql.CompileError){
    token, errorToken := lexical.Token()
    if errorToken != nil {
        return *fields, errorToken
    }

    if lexical.T_COMMA == token {

        token, errorToken = lexical.Token()
        if errorToken != nil {
            return *fields, errorToken
        }
        if token != lexical.T_LITERAL {
            return *fields, throwSyntaxError(lexical.T_LITERAL, token)
        }

        n := append(*fields, lexical.CurrentLexeme)
        fields = &n
        
        n, errorSyntax := g_table_params_rest(fields, count + 1)
        fields = &n
        if errorSyntax != nil {
            return *fields, errorSyntax
        }
    }

    return *fields, nil
}