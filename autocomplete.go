package main

import "github.com/k0kubun/pp"

func suggestTokens() {

}

func suggestColumns1(focus []rune, pos int) [][]rune {
	candidacies := [][]rune{
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

func suggestColumns2(focused string) [][]rune {
	focused = focused[:len(focused)-1]
	candidacies := [][]string{
		[]string{"hash", "date", "author", "author_email", "committer", "committer_email", "message", "full_message"},
		[]string{"name", "full_name", "type", "hash"},
		[]string{"name", "url", "push_url", "owner"},
		[]string{"name", "full_name", "hash"},
	}

	var suggests [][]rune
	for _, candidacy := range candidacies {
		s := containSlice(focused, candidacy)
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
	pp.Println(pos)
	//var suggest [][]rune
	ln := len(inputs)

	if ln == 1 {
		return [][]rune{[]rune("select")}
	} else if ln > 1 {
		focused := inputs[ln-2]
		focus := inputs[ln-1]

		// gitql> select remo[tab
		// In the case where the most recent input is "select"
		pp.Println(focused)
		if string(focused) == "select" {
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
		} else if containColumns(string(focused)) {
			if pos != 0 {
				return suggestColumns1(focus, pos)
			}
			return suggestColumns2(string(focused))
		}

	}

	//tokens := []string{"select", "from", "where", "order", "by", "or", "and", "limit", "in", "asc", "desc"}
	//tables := []string{"remotes", "tags", "branches", "commits", "refs"}

	return [][]rune{[]rune("select")}
}

func containSlice(focused string, candidacy []string) [][]rune {
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
	for i, val := range candidacy {
		if focused == val {
			return i, true
		}
	}
	return -1, false
}
