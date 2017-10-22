package runtime

import (
	"path/filepath"
	"testing"

	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/semantical"
)

func failTestIfError(err error, t *testing.T) {
	if err != nil {
		t.Error(err.Error())
	}
}

func getTableForQuery(query, directory string, t *testing.T) *TableData {
	parser.New(query)
	ast, errGit := parser.AST()
	failTestIfError(errGit, t)

	folder, errFile := filepath.Abs(directory)
	failTestIfError(errFile, t)
	ast.Path = &folder
	errGit = semantical.Analysis(ast)
	failTestIfError(errGit, t)

	builder = GetGitBuilder(ast.Path)
	visitor := new(RuntimeVisitor)
	err := visitor.Visit(ast)
	failTestIfError(err, t)

	tableData, err := walkCommits(ast, visitor)
	failTestIfError(err, t)
	return tableData
}
func TestSortOrdering(t *testing.T) {
	query := "select hash, date from commits order by date desc limit 3"
	tableData := getTableForQuery(query, "../", t)
	for i := 1; i < len(tableData.rows); i++ {
		if tableData.rows[i]["date"].(string) > tableData.rows[i-1]["date"].(string) {
			t.Errorf("Date not sored. row %d is bigger than row %d", i, i-1)
		}
	}

	queryWithoutDate := "select hash from commits order by date desc limit 3"
	tableDataNew := getTableForQuery(queryWithoutDate, "../", t)
	if len(tableData.rows) != len(tableDataNew.rows) {
		t.Error("Two queried returned different number of rows")
	}
	for i := 0; i < len(tableData.rows); i++ {
		if tableDataNew.rows[i]["hash"].(string) != tableData.rows[i]["hash"].(string) {
			t.Errorf("Data in row %d does not match on both tables", i)
		}
	}
}
