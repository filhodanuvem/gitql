package main

import "testing"

func TestSuggestTokensFromInputting(t *testing.T) {
	inputting := [][]rune{
		[]rune("sel"),
		[]rune("f"),
		[]rune("wh"),
		[]rune("ord"),
		[]rune("b"),
		[]rune("o"),
		[]rune("an"),
		[]rune("l"),
		[]rune("i"),
		[]rune("as"),
		[]rune("d"),
	}

	answer := [][]rune{
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

	for i, input := range inputting {
		result := suggestTokensFromInputting(input, len(input))
		if i == 5 {
			// inputting "o" test
			// should return [][]rune{[]rune("or"), []rune("order")}
			for _, got := range result {
				if !(string(got) == "or" || string(got) == "order") {
					t.Errorf("expected 'or' also 'order', got %s", string(got))
				}
			}
		} else {
			assertSuggests(t, string(answer[i]), string(result[0]))
		}
	}
}

func TestSuggestTablesFromInputting(t *testing.T) {
	inputting := [][]rune{
		[]rune("rem"),
		[]rune("t"),
		[]rune("br"),
		[]rune("c"),
		[]rune("ref"),
	}

	answer := [][]rune{
		[]rune("remotes"),
		[]rune("tags"),
		[]rune("branches"),
		[]rune("commits"),
		[]rune("refs"),
	}

	for i, input := range inputting {
		result := suggestTablesFromInputting(input, len(input))
		assertSuggests(t, string(answer[i]), string(result[0]))
	}
}

func TestSuggestColumnsFromInputting(t *testing.T) {
	inputting1 := [][]rune{
		[]rune("na"),
		[]rune("u"),
		[]rune("pu"),
		[]rune("own"),
		[]rune("h"),
		[]rune("da"),
		[]rune("me"),
		[]rune("t"),
	}

	answer1 := [][]rune{
		[]rune("name"),
		[]rune("url"),
		[]rune("push_url"),
		[]rune("owner"),
		[]rune("hash"),
		[]rune("date"),
		[]rune("message"),
		[]rune("type"),
	}

	for i, input := range inputting1 {
		result := suggestColumnsFromInputting(input, len(input))
		assertSuggests(t, string(answer1[i]), string(result[0]))
	}

	inputting2 := [][]rune{
		[]rune("ful"),
		[]rune("au"),
		[]rune("co"),
	}

	for i, input := range inputting2 {
		result := suggestColumnsFromInputting(input, len(input))
		for _, got := range result {
			if i == 0 && !(string(got) == "full_name" || string(got) == "full_message") {
				t.Errorf("expected 'full_name' also 'full_message', got %s", string(got))
			}

			if i == 1 && !(string(got) == "author" || string(got) == "author_email") {
				t.Errorf("expected 'author' also 'author_email', got %s", string(got))
			}

			if i == 2 && !(string(got) == "committer" || string(got) == "committer_email") {
				t.Errorf("expected 'committer' also 'committer_email', got %s", string(got))
			}
		}
	}
}

func TestSuggestColumnsFromLatest(t *testing.T) {
	answer := []string{
		"hash",
		"type",
		"date",
		"author",
		"author_email",
		"committer",
		"committer_email",
		"message",
		"full_message",
		"name",
		"full_name",
	}

	expected := createHashMap(answer)
	result := suggestColumnsFromLatest("hash,")
	for _, v := range result {
		if _, ok := expected[string(v)]; !ok {
			t.Errorf("expected 'hash', 'type', 'date', 'author', 'author_email', 'committer', 'committer_email', 'message', 'full_message', 'name', 'full_name' got %s", string(v))
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	words := [][]rune{
		[]rune("alpaca"),
		[]rune("alpaca"),
		[]rune("Code-Hex"),
		[]rune("Hello"),
		[]rune("Hello"),
		[]rune("World"),
	}

	removeDuplicates(&words)

	if len(words) != 4 {
		t.Error("Failed to remove duplicates")
	}

	if string(words[0]) != "alpaca" {
		t.Errorf("expected alpaca, got %s", string(words[0]))
	}

	if string(words[1]) != "Code-Hex" {
		t.Errorf("expected Code-Hex, got %s", string(words[1]))
	}

	if string(words[2]) != "Hello" {
		t.Errorf("expected Hello, got %s", string(words[2]))
	}

	if string(words[3]) != "World" {
		t.Errorf("expected World, got %s", string(words[3]))
	}
}

func TestIsContained(t *testing.T) {
	idx, ok := isContained("Code-Hex", []string{"cloudson", "luizperes", "Code-Hex"})
	if !ok {
		t.Error("Failed to invoke `isContained`")
	}

	if idx != 2 {
		t.Errorf("expected %d, got %d", 2, idx)
	}
}

func TestGetPartsFromSlice(t *testing.T) {
	gotSlice := getPartsFromSlice("luizperes", []string{"cloudson", "luizperes", "Code-Hex"})
	if len(gotSlice) != 2 {
		t.Error("Failed to invoke `getPartsFromSlice`")
	}

	if string(gotSlice[0]) != "cloudson" {
		t.Errorf("expected cloudson, got %s", string(gotSlice[0]))
	}

	if string(gotSlice[1]) != "Code-Hex" {
		t.Errorf("expected Code-Hex, got %s", string(gotSlice[1]))
	}

	notSlice := getPartsFromSlice("gopher", []string{"cloudson", "luizperes", "Code-Hex"})
	if len(notSlice) != 0 {
		t.Error("Failed to invoke `getPartsFromSlice`")
	}
}

func createHashMap(s []string) map[string]bool {
	h := make(map[string]bool, len(s))
	for _, key := range s {
		h[key] = true
	}
	return h
}

func assertSuggests(t *testing.T, expected string, got string) {
	if expected != got {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
