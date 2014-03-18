package parser

type NodeMain interface {
    Run()
}

type NodeProgram struct {
    child NodeMain
}

type NodeSelect struct {
    WildCard bool
    fields []string
    tables []string
}

func (s *NodeSelect) Run() {
    return 
}

type NodeEmpty struct {

}

func (e *NodeEmpty) Run() {
    return 
}