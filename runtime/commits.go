package runtime

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/utilities"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func walkCommits(n *parser.NodeProgram, visitor *RuntimeVisitor) (*TableData, error) {
	iter, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	s := n.Child.(*parser.NodeSelect)
	where := s.Where
	distinct := s.Distinct

	counter := 1
	fields := s.Fields
	if s.WildCard {
		fields = builder.possibleTables[s.Tables[0]]
	}
	resultFields := fields // These are the fields in output with wildcards expanded
	rows := make([]tableRow, 1)
	usingOrder := false
	if s.Order != nil && !s.Count {
		usingOrder = true
		// Check if the order by field is in the selected fields. If not, add them to selected fields list
		if !utilities.IsFieldPresentInArray(fields, s.Order.Field) {
			fields = append(fields, s.Order.Field)
		}
	}

	// holds the seen values so far. field -> (value -> wasSeen)
	seen := make(map[string]map[string]bool)
	iter.ForEach(func(commit *object.Commit) error {
		builder.setCommit(commit)
		boolRegister = true
		visitor.VisitExpr(where)

		if boolRegister {
			if !s.Count {
				newRow := make(tableRow)

				isNew := true
				for _, f := range fields {
					data := metadataCommit(f, commit)

					_, ok := seen[f]
					if !ok {
						seen[f] = make(map[string]bool)
					}

					isNew = !seen[f][data]

					newRow[f] = data
					seen[f][data] = true
				}

				if isNew || distinct == nil {
					counter = counter + 1
					rows = append(rows, newRow)
				}
			}
		}

		if !usingOrder && !s.Count && counter > s.Limit {
			return fmt.Errorf("limit") // stop iteration
		}

		return nil
	})

	if s.Count {
		newRow := make(tableRow)
		// counter was started from 1!
		newRow[COUNT_FIELD_NAME] = strconv.Itoa(counter - 1)
		counter = 2
		rows = append(rows, newRow)
	}

	rowsSliced := rows[len(rows)-counter+1:]
	rowsSliced, err = orderTable(rowsSliced, s.Order)
	if err != nil {
		return nil, err
	}

	if usingOrder && !s.Count && counter > s.Limit {
		counter = s.Limit
		rowsSliced = rowsSliced[0:counter]
	}

	tableData := new(TableData)
	tableData.rows = rowsSliced
	tableData.fields = resultFields

	return tableData, nil
}

func metadataCommit(identifier string, commit *object.Commit) string {
	key := ""
	for key, _ = range builder.tables {
		break
	}
	table := key
	err := builder.UseFieldFromTable(identifier, table)
	if err != nil {
		log.Fatalln(err)
	}

	switch identifier {
	case "hash":
		return commit.ID().String()
	case "author":
		return commit.Author.Name
	case "author_email":
		return commit.Author.Email
	case "committer":
		return commit.Committer.Name
	case "committer_email":
		return commit.Committer.Email
	case "date":
		//return object.Committer().When.Format()
		return commit.Author.When.Format(parser.Time_YMDHIS)
	case "full_message":
		return commit.Message
	case "message":
		// return first line of a commit message
		message := commit.Message
		r := []rune("\n")
		idx := strings.IndexRune(message, r[0])
		if idx != -1 {
			message = message[0:idx]
		}
		return message

	}
	log.Fatalf("Field %s not implemented yet \n", identifier)

	return ""
}
