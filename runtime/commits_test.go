package runtime

import (
	"path/filepath"
	"testing"

	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/semantical"
)

func TestSortOrdering(t *testing.T) {

	failTestIfError := func(err error) {
		if err != nil {
			t.Error(err.Error())
		}
	}
	folder, errFile := filepath.Abs("../")
	failTestIfError(errFile)

	query := "select hash, date from commits order by date desc limit 3"
	parser.New(query)
	ast, errGit := parser.AST()
	failTestIfError(errGit)

	ast.Path = &folder
	errGit = semantical.Analysis(ast)
	failTestIfError(errGit)

	builder = GetGitBuilder(ast.Path)
	visitor := new(RuntimeVisitor)
	err := visitor.Visit(ast)
	failTestIfError(err)

	tableData, err := walkCommits(ast, visitor)
	failTestIfError(err)
	for i := 1; i < len(tableData.rows); i++ {
		if tableData.rows[i]["date"].(string) > tableData.rows[i-1]["date"].(string) {
			t.Errorf("Date not sored. row %d is bigger than row %d", i, i-1)
		}
	}

	queryWithoutDate := "select hash from commits order by date desc limit 3"
	parser.New(queryWithoutDate)
	ast, errGit = parser.AST()
	failTestIfError(errGit)

	ast.Path = &folder
	errGit = semantical.Analysis(ast)
	failTestIfError(errGit)

	builder = GetGitBuilder(ast.Path)
	visitor = new(RuntimeVisitor)
	err = visitor.Visit(ast)
	failTestIfError(err)

	tableDataNew, err := walkCommits(ast, visitor)
	failTestIfError(err)
	if len(tableData.rows) != len(tableDataNew.rows) {
		t.Error("Two queried returned different number of rows")
	}

	for i := 0; i < len(tableData.rows); i++ {
		if tableDataNew.rows[i]["hash"].(string) != tableData.rows[i]["hash"].(string) {
			t.Errorf("Data in row %d does not match on both tables", i)
		}
	}
}
