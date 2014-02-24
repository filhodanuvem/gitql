package lexical 

var source string 
var currentPointer = 0

func New(s string) {
    source = s
}

func Token() uint8 {
    // return T_SELECT;

    var lexeme string
    for char := nextChar(); char != 0; char = nextChar() {
        lexeme = lexeme + string(char)
        if lexeme == " " {
            lexeme = ""
            continue
        }

        if lexeme == "select" {
            return T_SELECT
        }

        if lexeme == "*" {
            return T_WILD_CARD
        }

        if lexeme == "from" {
            return T_FROM 
        }
    }

    return 0
}

func nextChar() uint8 {
    defer func() {
        currentPointer = currentPointer + 1
    }()

    if currentPointer >= len(source) {
        return 0;
    } 

    return source[currentPointer]
}

func rewind() {
    currentPointer = 0;
}