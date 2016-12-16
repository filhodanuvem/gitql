package runtime

import (
	"fmt"
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/semantical"
	"github.com/crackcomm/go-clitable"
	"github.com/cloudson/git2go"
	"log"
	"encoding/json"
)

const (
	WALK_COMMITS    = 1
	WALK_REFERENCES = 2
	WALK_REMOTES    = 3
)

const (
	REFERENCE_TYPE_BRANCH = "branch"
	REFERENCE_TYPE_REMOTE = "remote"
	REFERENCE_TYPE_TAG    = "tag"
)

var repo *git.Repository
var builder *GitBuilder
var boolRegister bool

type tableRow map[string]interface{}
type proxyTable struct {
	table  string
	fields map[string]string
}

type GitBuilder struct {
	tables           map[string]string
	possibleTables   map[string][]string
	proxyTables      map[string]*proxyTable
	repo             *git.Repository
	currentWalkType  uint8
	currentCommit    *git.Commit
	currentReference *git.Reference
	currentRemote    *git.Remote
	walk             *git.RevWalk
}

type RuntimeError struct {
	code    uint8
	message string
}

type RuntimeVisitor struct {
	semantical.Visitor
}

type TableData struct {
	rows []tableRow
	fields []string
}

// =========================== Error

func (e *RuntimeError) Error() string {
	return e.message
}

func throwRuntimeError(message string, code uint8) *RuntimeError {
	e := new(RuntimeError)
	e.message = message
	e.code = code

	return e
}

// =========================== Runtime
func Run(n *parser.NodeProgram, typeFormat *string) {
	builder = GetGitBuilder(n.Path)
	visitor := new(RuntimeVisitor)
	err := visitor.Visit(n)
	if err != nil {
		log.Fatalln(err)
	}
	var tableData *TableData

	switch findWalkType(n) {
	case WALK_COMMITS:
		tableData = walkCommits(n, visitor)
		break
	case WALK_REFERENCES:
		tableData = walkReferences(n, visitor)
		break
	case WALK_REMOTES:
		tableData = walkRemotes(n, visitor)
		break
	}

	if *typeFormat == "json" {
		printJson(tableData)
	} else {
		printTable(tableData)
	}
}

func findWalkType(n *parser.NodeProgram) uint8 {
	s := n.Child.(*parser.NodeSelect)
	switch s.Tables[0] {
	case "commits":
		builder.currentWalkType = WALK_COMMITS
	case "remotes":
		builder.currentWalkType = WALK_REMOTES
	case "refs", "tags", "branches":
		builder.currentWalkType = WALK_REFERENCES
	}

	return builder.currentWalkType
}

func printTable(tableData *TableData) {
	table := clitable.New(tableData.fields)
	for _, r := range tableData.rows {
		table.AddRow(r)
	}
	table.Print()
}

func printJson(tableData *TableData) error {
	res, err := json.Marshal(tableData.rows)
	if err != nil {
		log.Fatalln(err)
		return throwRuntimeError(fmt.Sprintf("Json error:'%s'", err), 0)
	} else {
		fmt.Println(string(res))
		return nil
	}
}

func orderTable(rows []tableRow, order *parser.NodeOrder) []tableRow {
	if order == nil {
		return rows
	}
	// We will use parser.NodeGreater.Assertion(A, B) to know if
	// A > B and then switch their positions.
	// Unfortunaly, we will use bubble sort, that is O(n²)
	// @todo change to quick or other better sort.
	var orderer parser.NodeExpr
	if order.Asc {
		orderer = new(parser.NodeGreater)
	} else {
		orderer = new(parser.NodeSmaller)
	}

	field := order.Field
	key := ""
	for key, _ = range builder.tables {
		break
	}
	table := key
	err := builder.UseFieldFromTable(field, table)
	if err != nil {
		log.Fatalln(err)
	}

	for i, row := range rows {
		for j, rowWalk := range rows {
			if orderer.Assertion(fmt.Sprintf("%v", rowWalk[field]), fmt.Sprintf("%v", row[field])) {
				aux := rows[j]
				rows[j] = rows[i]
				rows[i] = aux
			}
		}
	}

	return rows
}

func metadata(identifier string) string {
	switch builder.currentWalkType {
	case WALK_COMMITS:
		return metadataCommit(identifier, builder.currentCommit)
	case WALK_REFERENCES:
		return metadataReference(identifier, builder.currentReference)
	case WALK_REMOTES:
		return metadataRemote(identifier, builder.currentRemote)
	}

	log.Fatalln("GOD!")

	return ""
}

// =================== GitBuilder

func GetGitBuilder(path *string) *GitBuilder {

	gb := new(GitBuilder)
	gb.tables = make(map[string]string)
	possibleTables := PossibleTables()
	gb.possibleTables = possibleTables

	proxyTables := map[string]*proxyTable{
		"tags":     proxyTableEntry("refs", map[string]string{"type": "tag"}),
		"branches": proxyTableEntry("refs", map[string]string{"type": "branch"}),
	}
	gb.proxyTables = proxyTables

	openRepository(path)

	gb.repo = repo

	return gb
}

func proxyTableEntry(t string, f map[string]string) *proxyTable {
	p := new(proxyTable)
	p.table = t
	p.fields = f

	return p
}

func openRepository(path *string) {
	_repo, err := git.OpenRepositoryExtended(*path)
	if err != nil {
		log.Fatalln(err)
	}
	repo = _repo
}

func (g *GitBuilder) setCommit(object *git.Commit) {
	g.currentCommit = object
}

func (g *GitBuilder) setReference(object *git.Reference) {
	g.currentReference = object
}

func (g *GitBuilder) setRemote(object *git.Remote) {
	g.currentRemote = object
}

func (g *GitBuilder) WithTable(tableName string, alias string) error {
	err := g.isValidTable(tableName)
	if err != nil {
		return err
	}

	if g.possibleTables[tableName] == nil {
		return throwRuntimeError(fmt.Sprintf("Table '%s' not found", tableName), 0)
	}

	if alias == "" {
		alias = tableName
	}

	g.tables[alias] = tableName

	return nil
}

func (g *GitBuilder) isProxyTable(tableName string) bool {
	_, isIn := g.proxyTables[tableName]

	return isIn
}

func  PossibleTables() (map[string][]string) {
	return map[string][]string{
		"commits": {
			"hash",
			"date",
			"author",
			"author_email",
			"committer",
			"committer_email",
			"message",
			"full_message",
		},
		"refs": {
			"name",
			"full_name",
			"type",
			"hash",
		},
		"remotes": {
			"name",
			"url",
			"push_url",
			"owner",
		},
		"tags": {
			"name",
			"full_name",
			"hash",
		},
		"branches": {
			"name",
			"full_name",
			"hash",
		},
	}
}

func (g *GitBuilder) isValidTable(tableName string) error {
	if _, isOk := g.possibleTables[tableName]; !isOk {
		return throwRuntimeError(fmt.Sprintf("Table '%s' not found", tableName), 0)
	}

	return nil
}

func (g *GitBuilder) UseFieldFromTable(field string, tableName string) error {
	err := g.isValidTable(tableName)
	if err != nil {
		return err
	}

	if field == "*" {
		return nil
	}

	table := g.possibleTables[tableName]
	for _, t := range table {
		if t == field {
			return nil
		}
	}

	return throwRuntimeError(fmt.Sprintf("Table '%s' has not field '%s'", tableName, field), 0)
}
