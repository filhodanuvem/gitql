package runtime

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloudson/git2go"
	"github.com/cloudson/gitql/parser"
)


func walkCommits(n *parser.NodeProgram, visitor *RuntimeVisitor) (*TableData, error){
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
	rows := make([]tableRow, s.Limit)
	usingOrder := false
	if s.Order != nil {
		usingOrder = true
	}
	fn := func(object *git.Commit) bool {
		builder.setCommit(object)
		boolRegister = true
		visitor.VisitExpr(where)
		if boolRegister {
			newRow := make(tableRow)
			for _, f := range fields {
				newRow[f] = metadataCommit(f, object)
			}
			rows = append(rows, newRow)

			counter = counter + 1
		}
		if !usingOrder && counter > s.Limit {
			return false
		}
		return true
	}

  err := builder.walk.Iterate(fn)
  if err != nil {
    fmt.Printf(err.Error())
  }
  rowsSliced := rows[len(rows)-counter+1:]
  rowsSliced, err = orderTable(rowsSliced, s.Order)
  if err != nil {
  	return nil, err
  }
  if usingOrder {
    if counter > s.Limit {
      counter = s.Limit
    }
    rowsSliced = rowsSliced[0:counter]
  }
  tableData := new(TableData)
  tableData.rows = rowsSliced
  tableData.fields = fields
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
