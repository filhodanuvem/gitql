package runtime

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cloudson/gitql/parser"
	"github.com/go-git/go-git/v5/plumbing"
)

func walkReferences(n *parser.NodeProgram, visitor *RuntimeVisitor) (*TableData, error) {
	s := n.Child.(*parser.NodeSelect)
	where := s.Where

	// @TODO make PR with Repository.WalkReference()
	iterator, err := builder.repo.References()
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

	iterator.ForEach(func(ref *plumbing.Reference) error {
		builder.setReference(ref)
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
					newRow[f] = metadataReference(f, ref)
				}
				rows = append(rows, newRow)
			}

			counter = counter + 1
			if !usingOrder && counter > s.Limit {
				return fmt.Errorf("limit")
			}
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

	if usingOrder && counter > s.Limit {
		counter = s.Limit
		rowsSliced = rowsSliced[0:counter]
	}

	tableData := new(TableData)
	tableData.rows = rowsSliced
	tableData.fields = fields

	return tableData, nil
}

func metadataReference(identifier string, ref *plumbing.Reference) string {
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
		return ref.Name().Short()
	case "full_name":
		return ref.Name().String()
	case "hash":
		target := ref.Hash()
		if target.IsZero() {
			return "NULL"
		}
		return target.String()
	case "type":
		if ref.Name().IsBranch() {
			return REFERENCE_TYPE_BRANCH
		}

		if ref.Name().IsTag() {
			return REFERENCE_TYPE_TAG
		}

		return "stash" // unknow
	}
	log.Fatalf("Field %s not implemented yet in reference\n", identifier)

	return ""
}
