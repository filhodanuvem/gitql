package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/runtime"
	"github.com/cloudson/gitql/semantical"
	"github.com/urfave/cli/v2"
)

const Version = "Gitql 2.1.0"

func main() {
	app := &cli.App{
		Name:        "gitql",
		Usage:       "A git query language",
		Version:     Version,
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "Enter to interactive mode",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   ".",
				Usage:   `The (optional) path to run gitql`,
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "table",
				Usage:   "The output type format {table|json}",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "show-tables",
				Aliases: []string{"s"},
				Usage:   "Show all tables",
				Action: func(c *cli.Context) error {
					fmt.Print("Tables: \n\n")

					for tableName, fields := range runtime.PossibleTables() {
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
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "The version of gitql",
				Action: func(c *cli.Context) error {
					fmt.Println(Version)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return cli.ShowAppHelp(c)
			}

			path, format := c.String("path"), c.String("format")
			if c.Bool("interactive") {
				return runPrompt(path, format)
			}

			return runQuery(c.Args().First(), path, format)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
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
