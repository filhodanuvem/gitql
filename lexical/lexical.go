package lexical 

var source string 
var currentPointer = 0

func New(s string) {
    source = s
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