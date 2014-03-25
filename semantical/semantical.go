package semantical

import (
    "github.com/cloudson/gitql/parser"
)

func analysis(ast *parser.NodeProgram) error {
    semantic := new(SemanticalVisitor)

    return semantic.Visit(ast)
}

