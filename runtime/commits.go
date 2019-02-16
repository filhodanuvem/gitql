package runtime

import (
	"fmt"
	"strconv"
	"log"
	"strings"

	"github.com/cloudson/git2go"
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/utilities"
)

func walkCommits(n *parser.NodeProgram, visitor *RuntimeVisitor) (*TableData, error) {
	builder.walk, _ = repo.Walk()
	builder.walk.PushHead()
	builder.walk.Sorting(git.SortTime)

	s := n.Child.(*parser.NodeSelect)
	where := s.Where

	counter := 1
	fields := s.Fields
	if s.WildCard {
		fields = builder.possibleTables[s.Tables[0]]
	}
	resultFields := fields // These are the fields in output with wildcards expanded
	rows := make([]tableRow, s.Limit)
	usingOrder := false
	if s.Order != nil && !s.Count {
		usingOrder = true
		// Check if the order by field is in the selected fields. If not, add them to selected fields list
		if !utilities.IsFieldPresentInArray(fields, s.Order.Field) {
			fields = append(fields, s.Order.Field)
		}
	}
	fn := func(object *git.Commit) bool {
		builder.setCommit(object)
		boolRegister = true
		visitor.VisitExpr(where)
		if boolRegister {
			if !s.Count {
				newRow := make(tableRow)
				for _, f := range fields {
					newRow[f] = metadataCommit(f, object)
				}
				rows = append(rows, newRow)
			}
			counter = counter + 1
		}
		if !usingOrder && !s.Count && counter > s.Limit {
			return false
		}
		return true
	}

	err := builder.walk.Iterate(fn)
	if err != nil {
		fmt.Printf(err.Error())
	}
	if s.Count {
		newRow := make(tableRow)
		// counter was started from 1!
		newRow[COUNT_FIELD_NAME] = strconv.Itoa(counter-1)
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

func metadataCommit(identifier string, object *git.Commit) string {
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
		return object.Id().String()
	case "author":
		return object.Author().Name
	case "author_email":
		return object.Author().Email
	case "committer":
		return object.Committer().Name
	case "committer_email":
		return object.Committer().Email
	case "date":
		return object.Committer().When.Format(parser.Time_YMDHIS)
	case "full_message":
		return object.Message()
	case "message":
		// return first line of a commit message
		message := object.Message()
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
