package parser

import (
    "fmt"
    "strconv"
    _"unicode"
    "github.com/cloudson/gitql/lexical"
)

var look_ahead uint8

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

func throwSyntaxError(expectedToken uint8, foundToken uint8) (error){
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
    node, err := g_program()
    program.child = node 

    return program, err
}

func g_program() (NodeMain, error) {
    token, err := lexical.Token()
    look_ahead = token
    
    if err != nil {
        return nil, err
    }

    s, err2 := gSelect()
    if s == nil {
        return new(NodeEmpty), err2
    }
    return s, nil
}

func gSelect() (*NodeSelect, error){
    if look_ahead != lexical.T_SELECT {
        return nil, throwSyntaxError(lexical.T_SELECT, look_ahead)
    }
    token, err := lexical.Token()
    look_ahead = token 
    if err != nil {
        return nil, err
    }
    s := new(NodeSelect)

    // PARAMETERS
    fields, err2 := gTableParams()
    if err2 != nil {
        return nil, err2
    }

    if len(fields) == 1 && fields[0] == "*" {
        s.WildCard = true   
    }
    s.fields = fields

    // TABLES
    tables , err4 := gTableNames()
    if err4 != nil {
        return nil, err4
    }
    s.tables = tables

    // LIMIT 
    var err5 error 
    s.limit, err5 = gLimit() 
    if err5 != nil {
        return nil, err5 
    }

    return s, nil
}

func gTableNames() ([]string, error){
    if look_ahead != lexical.T_FROM {
        return nil, throwSyntaxError(lexical.T_FROM, look_ahead)
    } 
    look_ahead, error := lexical.Token()
    if error != nil {
        return nil,error
    }
    if look_ahead != lexical.T_LITERAL {
        return nil, throwSyntaxError(lexical.T_LITERAL, look_ahead)
    }
    
    tables := make([]string, 1)
    tables[0] = lexical.CurrentLexeme
    
    return tables, nil
}

func gTableParams() ([]string, error){
    if look_ahead == lexical.T_WILD_CARD {
        token, err := lexical.Token()
        if err != nil {
            return nil,err
        }
        look_ahead = token
        return []string{"*"}, nil
    }
    var fields = []string{}
    if look_ahead == lexical.T_LITERAL {
        fields := append(fields, lexical.CurrentLexeme)
        token, err := lexical.Token()
        if err != nil {
            return nil,err
        }
        look_ahead = token
        fields, errorSyntax := gTableParamsRest(&fields, 1)

        return fields, errorSyntax
    }
    return nil, throwSyntaxError(lexical.T_LITERAL, look_ahead)
    
}

func gTableParamsRest(fields *[]string, count int) ([]string, error){
    if lexical.T_COMMA == look_ahead {
        var errorToken *lexical.TokenError
        look_ahead, errorToken = lexical.Token()
        if errorToken != nil {
            return *fields, errorToken
        }
        if look_ahead != lexical.T_LITERAL {
            return *fields, throwSyntaxError(lexical.T_LITERAL, look_ahead)
        }

        n := append(*fields, lexical.CurrentLexeme)
        fields = &n
        look_ahead, errorToken = lexical.Token()
        if errorToken != nil {
            return *fields, errorToken
        }        
        n, errorSyntax := gTableParamsRest(fields, count + 1)
        fields = &n
        if errorSyntax != nil {
            return *fields, errorSyntax
        }
    }

    return *fields, nil
}

func gLimit() (int, error) {
    if look_ahead != lexical.T_LIMIT {
        return 0, nil
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
    return number, nil  
}