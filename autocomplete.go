package main

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
		"distinct,",
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

func suggestQuery(inputs [][]rune, pos int) [][]rune {

	ln := len(inputs)

	if ln == 1 {
		// When nothing is input yet
		return [][]rune{[]rune("select")}
	}
	focused := string(inputs[ln-2])
	if focused == "select" {
		// gitql> select [tab
		// In the case where the most recent input is "select"
		return [][]rune{
			[]rune("distinct"),
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
		// gitql> select name, [tab
		// gitql> select committer, [tab
		// In the case where the most recent input is the column name and comma
		return suggestColumnsFromLatest(focused)
	} else if focused == "from" {
		// gitql> select * from [tab
		// In the case after inputted "from"
		return [][]rune{
			[]rune("tags"),
			[]rune("branches"),
			[]rune("commits"),
			[]rune("refs"),
		}
	} else if focused == "order" {
		return [][]rune{[]rune("by")}
	} else if focused == "where" || focused == "by" || focused == "or" || focused == "and" {
		// gitql> select name from commits where [tab
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
		case "branches", "tags":
			return [][]rune{
				[]rune("name"),
				[]rune("full_name"),
				[]rune("hash"),
			}
		}
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
	idx, ok := isContained(focused, candidacy)
	if ok {
		var suggests [][]rune
		for i, v := range candidacy {
			// Create slices other than what was focused
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
