package parser

type NodeMain interface {
    Run()
}

type NodeProgram struct {
    child NodeMain
}

type NodeSelect struct {
    WildCard bool
}