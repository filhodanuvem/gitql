package runtime

import (
	"fmt"
	"log"
	"os"
	"strings"

	"encoding/json"

	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/semantical"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/olekukonko/tablewriter"
)

const (
	WALK_COMMITS    = 1
	WALK_REFERENCES = 2
)

const (
	REFERENCE_TYPE_BRANCH = "branch"
	REFERENCE_TYPE_REMOTE = "remote"
	REFERENCE_TYPE_TAG    = "tag"
)

const (
	COUNT_FIELD_NAME = "count"
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
	tables          map[string]string
	possibleTables  map[string][]string
	proxyTables     map[string]*proxyTable
	repo            *git.Repository
	currentWalkType uint8
	currentCommit   *object.Commit

	currentReference *plumbing.Reference
	//walk             *object.RevWalk
}

type RuntimeError struct {
	code    uint8
	message string
}

type RuntimeVisitor struct {
	semantical.Visitor
}

type TableData struct {
	rows   []tableRow
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
func RunSelect(n *parser.NodeProgram, typeFormat *string) error {
	builder = GetGitBuilder(n.Path)
	visitor := new(RuntimeVisitor)
	err := visitor.Visit(n)
	if err != nil {
		return err
	}
	var tableData *TableData

	switch findWalkType(n) {
	case WALK_COMMITS:
		tableData, err = walkCommits(n, visitor)
		break
	case WALK_REFERENCES:
		tableData, err = walkReferences(n, visitor)
		break
	}

	if err != nil {
		return err
	}

	if *typeFormat == "json" {
		printJson(tableData)
	} else {
		printTable(tableData)
	}

	return nil
}

func RunShow(node *parser.NodeProgram) error {
	s := node.Child.(*parser.NodeShow)
	if s.Databases {
		builder = GetGitBuilder(node.Path)
		fmt.Print("Databases: \n\n")
		databases, err := PossibleDatabases()
		if err != nil {
			return err
		}
		for _, database := range databases {
			fmt.Println(database)
		}
		return nil
	} else if s.Tables {
		fmt.Print("Tables: \n\n")
		for tableName, fields := range PossibleTables() {
			fmt.Printf("%s\n\t", tableName)
			for i, field := range fields {
				comma := "."
				if i+1 < len(fields) {
					comma = ", "
				}
				fmt.Printf("%s%s", field, comma)
			}
			fmt.Println()
		}
	}
	return nil
}

func RunUse(node *parser.NodeProgram) error {
	builder = GetGitBuilder(node.Path)
	u := node.Child.(*parser.NodeUse)

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	s, err := w.Status()
	if err != nil {
		return err
	}

	if !s.IsClean() {
		return fmt.Errorf("worktree is not clean")
	}

	refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", u.Branch))
	cOp := &git.CheckoutOptions{
		Branch: refName,
		Create: false,
	}
	err = w.Checkout(cOp)
	if err != nil {
		// Try fetching branch from origin and then switching to it.
		// If it doesn't work, return the original error.
		remote, remoteErr := repo.Remote("origin")
		if remoteErr != nil {
			return err
		}
		remoteErr = remote.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{
				config.RefSpec(fmt.Sprintf("%s:%s", refName, refName)),
			},
		})
		if remoteErr != nil {
			return err
		}
		err = w.Checkout(cOp)
	}

	if err == nil {
		fmt.Println("switched to database", u.Branch)
	}

	return err
}

func findWalkType(n *parser.NodeProgram) uint8 {
	s := n.Child.(*parser.NodeSelect)
	switch s.Tables[0] {
	case "commits":
		builder.currentWalkType = WALK_COMMITS
	case "refs", "tags", "branches":
		builder.currentWalkType = WALK_REFERENCES
	}

	return builder.currentWalkType
}

func printTable(tableData *TableData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(tableData.fields)
	table.SetRowLine(true)
	for _, row := range tableData.rows {
		rowData := make([]string, len(tableData.fields))
		for i, field := range tableData.fields {
			rowData[i] = fmt.Sprintf("%v", row[field])
		}
		table.Append(rowData)
	}
	table.Render()
}

func printJson(tableData *TableData) error {
	res, err := json.Marshal(tableData.rows)
	if err != nil {
		return throwRuntimeError(fmt.Sprintf("Json error:'%s'", err), 0)
	} else {
		fmt.Println(string(res))
	}
	return nil
}

func orderTable(rows []tableRow, order *parser.NodeOrder) ([]tableRow, error) {
	if order == nil {
		return rows, nil
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
		return nil, err
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

	return rows, nil
}

func metadata(identifier string) string {
	switch builder.currentWalkType {
	case WALK_COMMITS:
		return metadataCommit(identifier, builder.currentCommit)
	case WALK_REFERENCES:
		return metadataReference(identifier, builder.currentReference)
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
	_repo, err := git.PlainOpen(*path)
	if err != nil {
		log.Fatalln(err)
	}
	repo = _repo
}

func (g *GitBuilder) setCommit(object *object.Commit) {
	g.currentCommit = object
}

func (g *GitBuilder) setReference(object *plumbing.Reference) {
	g.currentReference = object
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

func PossibleTables() map[string][]string {
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

func PossibleDatabases() ([]string, error) {
	// local branches
	iter, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	branches := make([]string, 0)
	iter.ForEach(func(r *plumbing.Reference) error {
		if r.Name().IsBranch() {
			branches = append(branches, r.Name().Short())
		}
		return nil
	})

	// remote branches
	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}
	refList, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	refPrefix := "refs/heads/"
	for _, ref := range refList {
		refName := ref.Name().String()
		if strings.HasPrefix(refName, refPrefix) {
			branches = append(branches, "remotes/origin/" + ref.Name().Short())
		}
	}

	return branches, nil
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
