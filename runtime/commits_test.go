package runtime

import (
	"fmt"
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
	findWalkType(ast)
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

func TestRowLimitsCount(t *testing.T) {
	query := "select hash, date from commits order by date desc limit 3"
	tableData := getTableForQuery(query, "../", t)

	if len(tableData.rows) > 3 {
		t.Error("Got more rows than the limit ")
	}
}

func TestWildcardFieldsCount(t *testing.T) {
	query := "select * from commits"
	table := getTableForQuery(query, "../", t)
	if len(table.fields) != 8 {
		t.Errorf("Commits has 8 fields. Output table got %d fields", len(table.fields))
	}
}

func TestSelectedFieldsCount(t *testing.T) {
	query := "select author, hash from commits"
	table := getTableForQuery(query, "../", t)
	if len(table.fields) != 2 {
		t.Errorf("Selected 2 fields. Output table got %d fields", len(table.fields))
	}
	if table.fields[0] != "author" || table.fields[1] != "hash" {
		t.Errorf("Selected 'author' and 'hash'. Got %v", table.fields)
	}
}

func TestNotEqualsInWhereLTGT(t *testing.T) {
	queryData := "select committer, hash from commits limit 1"
	table := getTableForQuery(queryData, "../", t)
	firstCommitter := table.rows[0]["committer"].(string)
	query := fmt.Sprintf("select committer, hash from commits where committer <> '%s' limit 1", firstCommitter)
	table = getTableForQuery(query, "../", t)
	if firstCommitter == table.rows[0]["committer"].(string) {
		t.Errorf("Still got the same committer as the first one. - %s", firstCommitter)
	}
}
func TestNotEqualsInWhere(t *testing.T) {
	queryData := "select committer, hash from commits limit 1"
	table := getTableForQuery(queryData, "../", t)
	firstCommitter := table.rows[0]["committer"].(string)
	query := fmt.Sprintf("select committer, hash from commits where committer != '%s' limit 1", firstCommitter)
	table = getTableForQuery(query, "../", t)
	if firstCommitter == table.rows[0]["committer"].(string) {
		t.Errorf("Still got the same committer as the first one. - %s", firstCommitter)
	}
}
