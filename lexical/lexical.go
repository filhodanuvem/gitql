package lexical 

import (
    "unicode"
    "strings"
    "fmt"
)

var source string 
var currentPointer int
var CurrentLexeme string

type TokenError struct {
    char int32
}

func (error *TokenError) Error() string {
    return fmt.Sprintf("Unexpected char '%s' with source '%s'", string(error.char), source)
}

func throwTokenError(char int32) (*TokenError) {
    error := new(TokenError)
    error.char = char 

    return error
}

var char int32
func New(s string) {
    source = s
    currentPointer = 0
    char = nextChar()
}

func Token() (uint8, *TokenError) {
    var lexeme string
    defer func() {
        CurrentLexeme = lexeme 
    }()
    state := S_START;
    for true {
        switch state {
            case S_START :
                if unicode.IsLetter(char) {
                    state = S_ID
                    break
                } else if unicode.IsNumber(char) {
                    state = S_NUMERIC
                } else {
                    lexeme = lexeme + string(char)
                    switch lexeme {
                        case "*" : 
                            state = S_WILD_CARD
                            break
                        case ",":
                            state = S_COMMA
                            break
                        case ";":
                            state = S_SEMICOLON
                            break
                        case ">":
                            state = S_GREATER
                            break
                        case "<":
                            state = S_SMALLER
                            break
                        case "=":
                            state = S_EQUAL
                            break
                        case "!" : 
                            state = S_NOT_EQUAL
                            break
                        case " ":
                            lexeme = ""
                            char = nextChar()
                            state = S_START
                            break
                        default: 
                            return 0, throwTokenError(char);
                    }
                }
                break;
            case S_ID: 
                for unicode.IsLetter(char) || unicode.IsNumber(char) || string(char) == "_" {
                    lexeme = lexeme + string(char)
                    char = nextChar()
                }
                return lexemeToToken(lexeme), nil
            case S_NUMERIC: 
                for unicode.IsNumber(char) {
                    lexeme = lexeme + string(char)
                    char = nextChar()
                }
                // @TODO insert to symbol table
                return T_NUMERIC, nil
            case S_WILD_CARD: 
                char = nextChar()
                return T_WILD_CARD, nil
            case S_COMMA: 
                char = nextChar()
                return T_COMMA, nil
            case S_SEMICOLON: 
                char = nextChar()
                return T_SEMICOLON, nil
            case S_GREATER: 
                char = nextChar()
                lexeme = string(char)
                if lexeme == "=" {
                    state = S_GREATER_OR_EQUAL
                    break
                }
                return T_GREATER, nil
            case S_GREATER_OR_EQUAL: 
                char = nextChar()
                return T_GREATER_OR_EQUAL, nil
            case S_SMALLER: 
                char = nextChar()
                lexeme = string(char)
                if lexeme == "=" {
                    state = S_SMALLER_OR_EQUAL
                    break
                }
                return T_SMALLER, nil
            case S_SMALLER_OR_EQUAL: 
                char = nextChar()
                return T_SMALLER_OR_EQUAL, nil
            case S_EQUAL: 
                char = nextChar()
                return T_EQUAL, nil
            case S_NOT_EQUAL: 
                char = nextChar()
                lexeme = string(char)
                if lexeme == "=" {
                    return T_NOT_EQUAL, nil
                }
                return 0, throwTokenError(char)
            default:
                break
        }
    }
    return T_EOF, throwTokenError(char)
}

func lexemeToToken(lexeme string) (uint8){
    switch strings.ToLower(lexeme) {
        case "select" : return T_SELECT
        case "from" : return T_FROM
        case "where" : return T_WHERE
        case "order" : return T_ORDER
        case "by" : return T_BY
        case "limit" : return T_LIMIT
    }
    // @TODO insert to symbol table
    return T_LITERAL
}

func nextChar() int32 {
    defer func() {
        currentPointer = currentPointer + 1
    }()

    if currentPointer >= len(source) {
        return T_EOF;
    } 

    return int32(source[currentPointer])
}

func rewind() {
    currentPointer = 0;
}