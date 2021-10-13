package lexical

const T_SELECT = 1
const T_DISTINCT = 29
const T_FROM = 2
const T_WHERE = 3
const T_ORDER = 4
const T_BY = 5
const T_LIMIT = 6
const T_DESC = 7
const T_WILD_CARD = 8
const T_COMMA = 9
const T_SEMICOLON = 10
const T_GREATER = 11
const T_SMALLER = 12
const T_GREATER_OR_EQUAL = 13
const T_SMALLER_OR_EQUAL = 14
const T_EQUAL = 15
const T_NOT_EQUAL = 16
const T_LITERAL = 17
const T_NUMERIC = 18
const T_AND = 19
const T_OR = 20
const T_ID = 21
const T_PARENTH_L = 22
const T_PARENTH_R = 23
const T_IN = 24
const T_ASC = 25
const T_LIKE = 26
const T_NOT = 27
const T_COUNT = 28
const T_EOF = 0
const T_FUCK = 66
const T_SHOW = 31
const T_TABLES = 32
const T_DATABASES = 33
const T_USE = 34

var tokenNamesByValue = map[uint8]string{
	T_SELECT:           "T_SELECT",
	T_DISTINCT:         "T_DISTINCT",
	T_FROM:             "T_FROM",
	T_WHERE:            "T_WHERE",
	T_ORDER:            "T_ORDER",
	T_BY:               "T_BY",
	T_LIMIT:            "T_LIMIT",
	T_DESC:             "T_DESC",
	T_WILD_CARD:        "T_WILD_CARD",
	T_COMMA:            "T_COMMA",
	T_SEMICOLON:        "T_SEMICOLON",
	T_GREATER:          "T_GREATER",
	T_SMALLER:          "T_SMALLER",
	T_GREATER_OR_EQUAL: "T_GREATER_OR_EQUAL",
	T_SMALLER_OR_EQUAL: "T_SMALLER_OR_EQUAL",
	T_EQUAL:            "T_EQUAL",
	T_NOT_EQUAL:        "T_NOT_EQUAL",
	T_LITERAL:          "T_LITERAL",
	T_NUMERIC:          "T_NUMERIC",
	T_ID:               "T_ID",
	T_PARENTH_L:        "T_PARENTH_L",
	T_PARENTH_R:        "T_PARENTH_R",
	T_IN:               "T_IN",
	T_EOF:              "T_EOF",
	T_ASC:              "T_ASC",
	T_NOT:              "T_NOT",
	T_COUNT:            "T_COUNT",
	T_SHOW:             "T_SHOW",
	T_TABLES:           "T_TABLES",
	T_DATABASES:        "T_DATABASES",
	T_USE:              "T_USE",
}

func TokenName(token uint8) string {
	return tokenNamesByValue[token]
}
