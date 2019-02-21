package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/runtime"
	"github.com/cloudson/gitql/semantical"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

const Version = "Gitql 1.6.0"

type Gitql struct {
	Path          string `short:"p" default:"."`
	Version       bool   `short:"v"`
	Isinteractive bool   `short:"i"`
	ShowTables    bool   `long:"show-tables"`
	TypeFormat    string `long:"type" default:"table"`
	Query         string
}

func (cmd *Gitql) Run() int {
	if err := unwrap(cmd.execute()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	return 0
}

func (cmd Gitql) execute() error {
	if err := cmd.parseCommandLine(); err != nil {
		return err
	}

	folder, err := filepath.Abs(cmd.Path)
	if err != nil {
		return err
	}

	if cmd.Isinteractive {
		return runPrompt(folder, cmd.TypeFormat)
	}

	return runQuery(cmd.Query, folder, cmd.TypeFormat)
}

func runPrompt(folder, typeFormat string) error {

	term, err := readline.NewEx(&readline.Config{
		Prompt:       "gitql> ",
		AutoComplete: readline.SegmentFunc(suggestQuery),
	})
	if err != nil {
		return err
	}
	defer term.Close()

	for {
		query, err := term.Readline()
		if err != nil {
			if err == io.EOF {
				break // Ctrl^D
			}
			return err
		}

		if query == "" {
			continue
		}

		if query == "exit" || query == "quit" {
			break
		}

		if err := runQuery(query, folder, typeFormat); err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}
	}

	return nil
}

func runQuery(query, folder, typeFormat string) error {
	parser.New(query)
	ast, err := parser.AST()
	if err != nil {
		return err
	}

	ast.Path = &folder
	if err := semantical.Analysis(ast); err != nil {
		return err
	}

	return runtime.Run(ast, &typeFormat)
}

func (cmd *Gitql) parseCommandLine() error {
	if err := cmd.parse(os.Args[1:]); err != nil {
		return err
	}

	if cmd.Version {
		return makeIgnoreErr(Version)
	}

	if cmd.ShowTables {
		return makeIgnoreErr(printTables())
	}

	return nil
}

func (cmd *Gitql) parse(argv []string) error {
	p := flags.NewParser(cmd, flags.PrintErrors)
	args, err := p.ParseArgs(argv)

	if err != nil {
		return err
	}

	if !cmd.Isinteractive && !cmd.Version && !cmd.ShowTables && len(args) == 0 {
		os.Stderr.Write(cmd.usage())
		return errors.New("invalid command line options")
	}

	cmd.Query = strings.Join(args, " ")

	return nil
}

func (cmd Gitql) usage() []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `Gitql - Git query language
Usage: gitql [flags] [args]

Flags:
  -i    Enter to interactive mode
  -p string
        The (optional) path to run gitql (default ".")
  --show-tables
        Show all tables
  --type string
	The output type format {table|json} (default "table")
  -v    The version of gitql
Arguments:
  sql: A query to run
`)

	return buf.Bytes()
}

func printTables() string {
	var buf bytes.Buffer

	buf.WriteString("Tables: \n\n")

	tables := runtime.PossibleTables()
	for tableName, fields := range tables {
		buf.WriteString(fmt.Sprintf("%s\n\t", tableName))
		for i, field := range fields {
			comma := "."
			if i+1 < len(fields) {
				comma = ", "
			}
			buf.WriteString(fmt.Sprintf("%s%s", field, comma))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func makeIgnoreErr(str string) error {
	return ignore{err: errors.New(str)}
}

// Ignore error
type ignore struct {
	err error
}

type cause interface {
	Cause() error
}

func (i ignore) Error() string {
	return i.err.Error()
}

func (i ignore) Cause() error {
	return i.err
}

// get important message from wrapped error message
func unwrap(err error) error {
	for e := err; e != nil; {
		switch e.(type) {
		case ignore:
			fmt.Println(e.Error())
			return nil
		case cause:
			e = e.(cause).Cause()
		default:
			return e
		}
	}

	return nil
}
