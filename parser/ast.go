package parser

type NodeMain interface {
    Run()
}

type NodeProgram struct {
    child NodeMain
}