package main

import (
	"strings"
	"testing"
)

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

func TestSuggestQuery(t *testing.T) {
	// gitql> [tab
	// expected: select
	pattern1 := [][]rune{
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern1, []string{"select"})

	// gitql> select [tab
	// expected: *, name, url,  push_url, owner, full_name, hash, date, author,
	// author_email, committer, committer_email, message, full_message, type
	pattern2 := [][]rune{
		[]rune("select"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern2, []string{
		"*",
		"name",
		"url",
		"push_url",
		"owner",
		"full_name",
		"hash",
		"date",
		"author",
		"author_email",
		"committer",
		"committer_email",
		"message",
		"full_message",
		"type",
	})

	// gitql> select name [tab
	// expected: select, from, where, order, by, or, and, limit, in, asc, desc
	pattern3 := [][]rune{
		[]rune("select"),
		[]rune("name"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern3, []string{
		"select",
		"from",
		"where",
		"order",
		"by",
		"or",
		"and",
		"limit",
		"in",
		"asc",
		"desc",
	})
	// gitql> select name, [tab
	// expected: full_name, type, hash, url, push_url, owner
	pattern4 := [][]rune{
		[]rune("select"),
		[]rune("name,"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern4, []string{
		"full_name",
		"type",
		"hash",
		"url",
		"push_url",
		"owner",
	})

	// gitql> select * from [tab
	// expected: tags, branches, commits, refs
	pattern5 := [][]rune{
		[]rune("select"),
		[]rune("*"),
		[]rune("from"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern5, []string{
		"tags",
		"branches",
		"commits",
		"refs",
	})

	// gitql> select name from refs where [tab
	// expected: name, url, push_url, owner
	pattern6 := [][]rune{
		[]rune("select"),
		[]rune("name"),
		[]rune("from"),
		[]rune("refs"),
		[]rune("where"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern6, []string{"name", "full_name", "type", "hash"})

	// gitql> select committer from commits where committer = "K" and [tab
	// expected: hash, date, author, author_email, committer, committer_email, message, full_message
	pattern7 := [][]rune{
		[]rune("select"),
		[]rune("committer"),
		[]rune("from"),
		[]rune("commits"),
		[]rune("where"),
		[]rune("committer"),
		[]rune("="),
		[]rune(`"K"`),
		[]rune("and"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern7, []string{
		"hash",
		"date",
		"author",
		"author_email",
		"committer",
		"committer_email",
		"message",
		"full_message",
	})

	// gitql> select committer from commits where committer = "k" order [tab
	// expected: by
	pattern8 := [][]rune{
		[]rune("select"),
		[]rune("committer"),
		[]rune("from"),
		[]rune("commits"),
		[]rune("where"),
		[]rune("committer"),
		[]rune("="),
		[]rune(`"K"`),
		[]rune("order"),
		[]rune(""),
	}
	assertSuggestsQuery(t, pattern8, []string{"by"})
}

// tiny tools
func assertSuggestsQuery(t *testing.T, inputs [][]rune, expected []string) {
	result := suggestQuery(inputs, len(inputs[len(inputs)-1]))
	expectedHash := createHashMap(expected)

	for _, v := range result {
		_, ok := expectedHash[string(v)]
		if !ok {
			t.Errorf("expected: (%s), got: %s", strings.Join(expected, ", "), string(v))
			break
		}
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
