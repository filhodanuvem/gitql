package runtime

import (
	"log"
	"strconv"

	git "github.com/libgit2/git2go"
	"github.com/cloudson/gitql/parser"
)

func walkReferences(n *parser.NodeProgram, visitor *RuntimeVisitor) (*TableData, error) {
	s := n.Child.(*parser.NodeSelect)
	where := s.Where

	// @TODO make PR with Repository.WalkReference()
	iterator, err := builder.repo.NewReferenceIterator()
	if err != nil {
		return nil, err
	}
	counter := 1
	fields := s.Fields
	if s.WildCard {
		fields = builder.possibleTables[s.Tables[0]]
	}
	rows := make([]tableRow, s.Limit)
	usingOrder := false
	if s.Order != nil && !s.Count {
		usingOrder = true
	}
	for object, inTheEnd := iterator.Next(); inTheEnd == nil; object, inTheEnd = iterator.Next() {

		builder.setReference(object)
		boolRegister = true
		visitor.VisitExpr(where)
		if boolRegister {
			fields := s.Fields
			if s.WildCard {
				fields = builder.possibleTables[s.Tables[0]]
			}
			if !s.Count {
				newRow := make(tableRow)
				for _, f := range fields {
					newRow[f] = metadataReference(f, object)
				}
				rows = append(rows, newRow)
			}
			counter = counter + 1
			if !usingOrder && counter > s.Limit {
				break
			}
		}
	}
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
	if usingOrder && counter > s.Limit {
		counter = s.Limit
		rowsSliced = rowsSliced[0:counter]
	}
	tableData := new(TableData)
	tableData.rows = rowsSliced
	tableData.fields = fields
	return tableData, nil
}

func metadataReference(identifier string, object *git.Reference) string {
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
	case "name":
		return object.Shorthand()
	case "full_name":
		return object.Name()
	case "hash":
		target := object.Target()
		if target == nil {
			return "NULL"
		}
		return target.String()
	case "type":
		if object.IsBranch() {
			return REFERENCE_TYPE_BRANCH
		}

		if object.IsRemote() {
			return REFERENCE_TYPE_REMOTE
		}

		if object.IsTag() {
			return REFERENCE_TYPE_TAG
		}

		return "stash" // unknow
	}
	log.Fatalf("Field %s not implemented yet in reference\n", identifier)

	return ""
}
