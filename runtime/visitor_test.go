package runtime

import (
	"path/filepath"
	"testing"

	"github.com/navigaid/gitql/parser"
	"github.com/navigaid/gitql/semantical"
)

func TestTestAllFieldsInExprBranches(t *testing.T) {
	query := "select * from branches where name = 'something' and somthing > 'name'"
	err := parseAndVisitQuery(query, "../", t)
	if err == nil {
		t.Error("Expected error, received none")
	}
}

func TestTestAllFieldsInExprBranchesWithCount(t *testing.T) {
	query := "select count(*) from branches where name = 'something' and somthing > 'name'"
	err := parseAndVisitQuery(query, "../", t)
	if err == nil {
		t.Error("Expected error, received none")
	}
}

func TestTestAllFieldsInExprRefs(t *testing.T) {
	query := "select * from refs where name = 'something' or type = 'asdfasdfsd'"
	err := parseAndVisitQuery(query, "../", t)
	if err != nil {
		t.Errorf("Unexpedted error %s", err)
	}
}

func TestTestAllFieldsInExprTags(t *testing.T) {
	query := "select * from tags where type = 'blah'"
	err := parseAndVisitQuery(query, "../", t)
	if err == nil {
		t.Errorf("Unexpedted error %s", err)
	}
}

func parseAndVisitQuery(query, dir string, t *testing.T) error {
	parser.New(query)
	ast, errGit := parser.AST()
	failTestIfError(errGit, t)

	folder, errFile := filepath.Abs(dir)
	failTestIfError(errFile, t)
	ast.Path = &folder
	errGit = semantical.Analysis(ast)
	failTestIfError(errGit, t)

	builder = GetGitBuilder(ast.Path)
	visitor := new(RuntimeVisitor)
	err := visitor.Visit(ast)
	return err
}
