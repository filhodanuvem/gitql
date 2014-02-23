package lexical

import "testing"

func setUp() {
    rewind()
}

func TestGetNextChar (t *testing.T) {
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

func TestEndOfFile (t *testing.T) {
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


func assertChar(t *testing.T, expected uint8, found uint8) {
    if found != expected {
        t.Errorf("Char '%s' is not '%s'", string(found), string(expected));
    }
}