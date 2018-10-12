package lexical

import (
	"fmt"
	"strings"
	"unicode"
)

var source string
var currentPointer int
var CurrentLexeme string

type TokenError struct {
	char int32
}

func (error *TokenError) Error() string {
	character := string(error.char)
	if char == T_EOF {
		character = "EOF"
	}
	return fmt.Sprintf("Unexpected char '%s' with source '%s'", character, source)
}

func throwTokenError(char int32) *TokenError {
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
	state := S_START
	for true {
		switch state {
		case S_START:
			if unicode.IsLetter(char) {
				state = S_ID
				break
			} else if unicode.IsNumber(char) {
				state = S_NUMERIC
			} else {
				lexeme = lexeme + string(char)
				switch lexeme {
				case "*":
					state = S_WILD_CARD
					break
				case "(":
					state = S_PARENTH_L
					break
				case ")":
					state = S_PARENTH_R
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
				case "!":
					state = S_NOT_EQUAL
					break
				case "'":
					lexeme = ""
					char = nextChar()
					state = S_LITERAL
					break
				case "\"":
					lexeme = ""
					char = nextChar()
					state = S_LITERAL_2
					break
				case " ":
					lexeme = ""
					char = nextChar()
					state = S_START
					break
				default:
					if char == T_EOF {
						return T_EOF, nil
					}
					return T_FUCK, throwTokenError(char)
				}
			}
			break
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
			} else if lexeme == ">" {
				char = nextChar()
				return T_NOT_EQUAL, nil
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
				char = nextChar()
				return T_NOT_EQUAL, nil
			}
			return 0, throwTokenError(char)
		case S_LITERAL:
			for string(char) != "'" && char != T_EOF {
				lexeme = lexeme + string(char)
				char = nextChar()
			}
			if char == T_EOF {
				return 0, throwTokenError(char)
			}
			char = nextChar()
			return T_LITERAL, nil
		case S_LITERAL_2:
			for string(char) != "\"" && char != T_EOF {
				lexeme = lexeme + string(char)
				char = nextChar()
			}
			if char == T_EOF {
				return 0, throwTokenError(char)
			}
			char = nextChar()
			return T_LITERAL, nil
		case S_PARENTH_L:
			char = nextChar()
			return T_PARENTH_L, nil
		case S_PARENTH_R:
			char = nextChar()
			return T_PARENTH_R, nil
		default:
			state = S_START
		}
	}
	return T_EOF, throwTokenError(char)
}

func lexemeToToken(lexeme string) uint8 {
	switch strings.ToLower(lexeme) {
	case L_SELECT:
		return T_SELECT
	case L_FROM:
		return T_FROM
	case L_WHERE:
		return T_WHERE
	case L_ORDER:
		return T_ORDER
	case L_BY:
		return T_BY
	case L_OR:
		return T_OR
	case L_AND:
		return T_AND
	case L_LIMIT:
		return T_LIMIT
	case L_IN:
		return T_IN
	case L_ASC:
		return T_ASC
	case L_DESC:
		return T_DESC
	case L_LIKE:
		return T_LIKE
	case L_NOT:
		return T_NOT
	}
	return T_ID
}

func nextChar() int32 {
	defer func() {
		currentPointer = currentPointer + 1
	}()

	if currentPointer >= len(source) {
		return T_EOF
	}

	return int32(source[currentPointer])
}

func rewind() {
	currentPointer = 0
}
