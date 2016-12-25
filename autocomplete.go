package main

func suggestTokensFromInputting(focus []rune, pos int) [][]rune {
	return suggestInputting(focus, pos, [][]rune{
		[]rune("select"),
		[]rune("from"),
		[]rune("where"),
		[]rune("order"),
		[]rune("by"),
		[]rune("or"),
		[]rune("and"),
		[]rune("limit"),
		[]rune("in"),
		[]rune("asc"),
		[]rune("desc"),
	})
}

func suggestTablesFromInputting(focus []rune, pos int) [][]rune {
	return suggestInputting(focus, pos, [][]rune{
		[]rune("remotes"),
		[]rune("tags"),
		[]rune("branches"),
		[]rune("commits"),
		[]rune("refs"),
	})
}

func suggestColumnsFromInputting(focus []rune, pos int) [][]rune {
	return suggestInputting(focus, pos, [][]rune{
		[]rune("name"),
		[]rune("url"),
		[]rune("push_url"),
		[]rune("owner"),
		[]rune("full_name"),
		[]rune("hash"),
		[]rune("date"),
		[]rune("author"),
		[]rune("author_email"),
		[]rune("committer"),
		[]rune("committer_email"),
		[]rune("message"),
		[]rune("full_message"),
		[]rune("type"),
	})
}

func suggestInputting(focus []rune, pos int, candidacies [][]rune) [][]rune {
	if len(focus) == 0 {
		return candidacies
	}

	var suggests [][]rune
	for _, candidacy := range candidacies {
		if len(candidacy) < pos {
			continue
		}

		var v rune
		for i := 0; i < pos; i++ {
			v |= focus[i] ^ candidacy[i]
		}
		if v == rune(0) {
			suggests = append(suggests, candidacy)
		}
	}

	return suggests
}

func suggestColumnsFromLatest(focused string) [][]rune {
	return suggestLatest(focused[:len(focused)-1], [][]string{
		[]string{"hash", "date", "author", "author_email", "committer", "committer_email", "message", "full_message"},
		[]string{"name", "full_name", "type", "hash"},
		[]string{"name", "url", "push_url", "owner"},
		[]string{"name", "full_name", "hash"},
	})
}

// Creates a candidate from the input previous character string.
func suggestLatest(focused string, candidacies [][]string) [][]rune {
	var suggests [][]rune
	for _, candidacy := range candidacies {
		s := getPartsFromSlice(focused, candidacy)
		if s != nil {
			suggests = append(suggests, s...)
		}
	}
	removeDuplicates(&suggests)

	return suggests
}

func containColumns(focused string) bool {
	_, ok := isContained(focused, []string{
		"select", // gitql> select [tab
		"name,",
		"url,",
		"push_url,",
		"owner,",
		"full_name,",
		"hash,",
		"date,",
		"author,",
		"author_email,",
		"committer,",
		"committer_email,",
		"message,",
		"full_message,",
		"type,",
	})
	return ok
}

// Remove duplicates elements in slice.
func removeDuplicates(s *[][]rune) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *s {
		key := string(x)
		if !found[key] {
			found[key] = true
			(*s)[j] = (*s)[i]
			j++
		}
	}
	*s = (*s)[:j]
}

func suggestCommands(inputs [][]rune, pos int) [][]rune {

	ln := len(inputs)

	if ln == 1 {
		// When nothing is input yet
		return [][]rune{[]rune("select")}
	}
	focused := string(inputs[ln-2])
	focus := inputs[ln-1]
	if focused == "select" {
		// gitql> select [tab
		// In the case where the most recent input is "select"
		return [][]rune{
			[]rune("*"),
			[]rune("name"),
			[]rune("url"),
			[]rune("push_url"),
			[]rune("owner"),
			[]rune("full_name"),
			[]rune("hash"),
			[]rune("date"),
			[]rune("author"),
			[]rune("author_email"),
			[]rune("committer"),
			[]rune("committer_email"),
			[]rune("message"),
			[]rune("full_message"),
			[]rune("type"),
		}
	} else if containColumns(focused) {
		if pos > 0 {
			// gitql> select na[tab
			// gitql> select commi[tab
			// In the case is inputting column
			return suggestColumnsFromInputting(focus, pos)
		}

		// gitql> select name, [tab
		// gitql> select committer, [tab
		// In the case where the most recent input is the column name and comma
		return suggestColumnsFromLatest(focused)
	} else if focused == "from" {
		// gitql> select * from re[tab
		// gitql> select * from bran[tab
		// In the case is inputting table name
		if pos > 0 {
			return suggestTablesFromInputting(focus, pos)
		}

		// gitql> select * from [tab
		// In the case after inputted "from"
		return [][]rune{
			[]rune("remotes"),
			[]rune("tags"),
			[]rune("branches"),
			[]rune("commits"),
			[]rune("refs"),
		}
	} else if focused == "order" {
		return [][]rune{[]rune("by")}
	} else if focused == "where" || focused == "by" || focused == "or" || focused == "and" {
		if pos > 0 {
			// gitql> select name from remotes where na[tab
			// gitql> select * from commits where committer = "K" order by com[tab
			// gitql> select * from commits where committer = "K" or com[tab
			// In the case is inputting column inputted after "where", "by", "and", "or"
			return suggestColumnsFromInputting(focus, pos)
		}

		// gitql> select name from remotes where [tab
		// gitql> select * from commits where committer = "K" order by [tab
		// gitql> select * from commits where committer = "K" and [tab
		// In the case is inputted after "where", "by", "and", "or"
		var table string
		for i := 0; i < len(inputs); i++ {
			if string(inputs[i]) == "from" {
				i++
				table = string(inputs[i])
			}
		}

		switch table {
		case "commits":
			return [][]rune{
				[]rune("hash"),
				[]rune("date"),
				[]rune("author"),
				[]rune("author_email"),
				[]rune("committer"),
				[]rune("committer_email"),
				[]rune("message"),
				[]rune("full_message"),
			}
		case "refs":
			return [][]rune{
				[]rune("name"),
				[]rune("full_name"),
				[]rune("type"),
				[]rune("hash"),
			}
		case "remotes":
			return [][]rune{
				[]rune("name"),
				[]rune("url"),
				[]rune("push_url"),
				[]rune("owner"),
			}
		case "tags":
			return [][]rune{
				[]rune("name"),
				[]rune("full_name"),
				[]rune("hash"),
			}
		case "branches":
			return [][]rune{
				[]rune("name"),
				[]rune("full_name"),
				[]rune("hash"),
			}
		}
	}

	// Other case
	// gitql> select * fr[tab
	// gitql> select date, message from commits wh[tab
	if pos > 0 {
		return suggestTokensFromInputting(focus, pos)
	}

	return [][]rune{
		[]rune("select"),
		[]rune("from"),
		[]rune("where"),
		[]rune("order"),
		[]rune("by"),
		[]rune("or"),
		[]rune("and"),
		[]rune("limit"),
		[]rune("in"),
		[]rune("asc"),
		[]rune("desc"),
	}
}

func getPartsFromSlice(focused string, candidacy []string) [][]rune {
	idx, isContained := isContained(focused, candidacy)
	if isContained {
		var suggests [][]rune
		for i, v := range candidacy {
			if i != idx {
				suggests = append(suggests, []rune(v))
			}
		}
		return suggests
	}
	return nil
}

func isContained(focused string, candidacy []string) (int, bool) {
	for idx, val := range candidacy {
		if focused == val {
			return idx, true
		}
	}
	return -1, false
}
