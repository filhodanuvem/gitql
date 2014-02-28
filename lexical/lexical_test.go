package lexical

import "testing"

func setUp() {
    rewind()
}

func TestGetNextChar(t *testing.T) {
    setUp()
    source = "gopher"

    var expected int32
    expected = 'g'
    char := nextChar()
    assertChar(t, expected, char)
    
    expected = 'o'
    char = nextChar()
    assertChar(t, expected, char)

    expected = 'p'
    char = nextChar()
    assertChar(t, expected, char)
}

func TestEndOfFile(t *testing.T) {
    setUp()
    source = "go"

    var expected int32
    expected = 'g'
    char := nextChar()
    assertChar(t, expected, char)
    
    expected = 'o'
    char = nextChar()
    assertChar(t, expected, char)

    expected = 0
    char = nextChar()
    assertChar(t, expected, char)
}

func TestRecognizeAnToken(t *testing.T) {
    setUp()
    source = ";"
    char = nextChar()

    var token uint8
    token, _ = Token()

    assertToken(t, token, T_SEMICOLON)
}

func TestRecognizeASequenceOfTokens(t *testing.T) {
    setUp()
    source = "*,>"
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_WILD_CARD)

    token, _ = Token()
    assertToken(t, token, T_COMMA)

    token, _ = Token()
    assertToken(t, token, T_GREATER)
}

func TestRecognizeTokensWithLexemesOfTwoChars(t *testing.T) {
    setUp()
    source = ">= <="
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_GREATER_OR_EQUAL)

    token, _ = Token()
    assertToken(t, token, T_SMALLER_OR_EQUAL)
}

func TestRecognizeTokensWithSourceManySpaced(t *testing.T) {
    setUp()
    source = "=    <    >=   !="
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_EQUAL)

    token, _ = Token()
    assertToken(t, token, T_SMALLER)    

    token, _ = Token()
    assertToken(t, token, T_GREATER_OR_EQUAL)

    token, _ = Token()
    assertToken(t, token, T_NOT_EQUAL)
}

func TestErrorUnrecognizeChar(t* testing.T) {
    setUp()
    source = "!"
    char = nextChar()

    _, error := Token()
    if (error == nil) {
        t.Errorf("Expected error with char '!' ")
    }
}

func TestReservedWords(t* testing.T) {
    setUp()
    source = "SELECT from WHEre"
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_SELECT)

    token, _ = Token()
    assertToken(t, token, T_FROM)

    token, _ = Token()
    assertToken(t, token, T_WHERE)
}

func TestNotReservedWords(t *testing.T) {
    setUp()

    source = "users commits"
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_LITERAL)

    token, _ = Token()
    assertToken(t, token, T_LITERAL)

}

func TestNumbers(t *testing.T) {
    setUp()

    source = "314 555"
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_NUMERIC)
}

func TestCurrentLexeme(t *testing.T) {
    setUp()
    source = "select * users"
    char = nextChar()

    var token uint8

    token, _ = Token()
    assertToken(t, token, T_SELECT)

    if (CurrentLexeme != "select") {
        t.Errorf("%s is not select", CurrentLexeme)
    }

    token, _ = Token()
    assertToken(t, token, T_WILD_CARD)

    if (CurrentLexeme != "*") {
        t.Errorf("%s is not *", CurrentLexeme)
    }


    token, _ = Token()
    assertToken(t, token, T_LITERAL)

    if (CurrentLexeme != "users") {
        t.Errorf("%s is not users", CurrentLexeme)
    }
}

func TestRepetitiveTokens(t *testing.T) {
    setUp()

    source = "select name, age from users"
    char = nextChar()
    
    var token uint8

    tokens := []uint8{T_SELECT, T_LITERAL, T_COMMA, T_LITERAL, T_FROM, T_LITERAL}
    for i := range tokens {
        token, _ = Token()
        assertToken(t, token, tokens[i])
    }
}

func assertToken(t *testing.T, expected uint8, found uint8) {
    if (expected != found) {
        t.Errorf("Token %s is not %s, lexeme: %s", TokenName(found), TokenName(expected), CurrentLexeme)
    }
}

func assertChar(t *testing.T, expected int32, found int32) {
    if found != expected {
        t.Errorf("Char '%s' is not '%s'", string(found), string(expected));
    }
}