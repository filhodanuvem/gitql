package lexical

import "testing"

func setUp() {
    rewind()
}

func TestGetNextChar(t *testing.T) {
    setUp()
    source = "gopher"

    var expected uint8
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

    var expected uint8
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
    source = "select"

    var token uint8
    token = Token()

    assertToken(t, token, T_SELECT)
}

func TestRecognizeASequenceOfTokens(t *testing.T) {
    setUp()
    source = "select * from"

    var token uint8
    token = Token()
    assertToken(t, token, T_SELECT)

    token = Token()
    assertToken(t, token, T_WILD_CARD)

    token = Token()
    assertToken(t, token, T_FROM)

}

func assertToken(t *testing.T, expected uint8, found uint8) {
    if (expected != found) {
        t.Errorf("Token %d is not %d", found, expected)
    }
}

func assertChar(t *testing.T, expected uint8, found uint8) {
    if found != expected {
        t.Errorf("Char '%s' is not '%s'", string(found), string(expected));
    }
}