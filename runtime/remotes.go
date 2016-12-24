package runtime

import (
	"log"

	"github.com/cloudson/git2go"
	"github.com/cloudson/gitql/parser"
)

func walkRemotes(n *parser.NodeProgram, visitor *RuntimeVisitor) (*TableData, error){
  s := n.Child.(*parser.NodeSelect)
  where := s.Where

	remoteNames, err := builder.repo.ListRemotes()
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
	if s.Order != nil {
		usingOrder = true
	}
	for _, remoteName := range remoteNames {
		object, errRemote := builder.repo.LoadRemote(remoteName)
		if errRemote != nil {
			log.Fatalln(errRemote)
		}

		builder.setRemote(object)
		boolRegister = true
		visitor.VisitExpr(where)
		if boolRegister {
			newRow := make(map[string]interface{})
			for _, f := range fields {
				newRow[f] = metadataRemote(f, object)
			}
			rows = append(rows, newRow)

      counter = counter + 1
      if !usingOrder && counter > s.Limit {
      	break
    	}
    }
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

func metadataRemote(identifier string, object *git.Remote) string {
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
		return object.Name()
	case "url":
		return object.Url()
	case "push_url":
		return object.PushUrl()
	case "owner":
		repo := object.Owner()
		r := &repo
		return r.Path()
	}

	log.Fatalf("Field %s not implemented yet \n", identifier)

	return ""
}
