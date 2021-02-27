package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestNoMatch(t *testing.T) {
	_, err := GetAdjacent("", strings.NewReader("<ta<ta<bl>"))

	if err == nil {
		t.Fail()
	}
}

func TestQuery(t *testing.T) {
	query := "X"
	expect := "y"
	content := fmt.Sprintf(
		"<table><tr><td>%s</td><td>%s</td></tr></table>",
		query,
		expect,
	)

	result, err := GetAdjacent(query, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestBackwardsQuery(t *testing.T) {
	query := "X"
	expect := "y"
	content := fmt.Sprintf(
		"<table><tr><td>%s</td><td>%s</td></tr></table>",
		expect,
		query,
	)

	result, err := GetAdjacent(query, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestTwo(t *testing.T) {
	query := "X"
	expect := "y"
	content := fmt.Sprintf(
		"<table><tr><td>a</td><td>b</td></tr><tr><td>%s</td><td>%s</td></tr></table>",
		expect,
		query,
	)

	result, err := GetAdjacent(query, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}
func TestDiacritics(t *testing.T) {
	query := "&#233;"
	expect := "b"
	content := fmt.Sprintf(
		"<table><tr><td>%s</td><td>%s</td></tr></table>",
		expect,
		query,
	)

	result, err := GetAdjacent("é", strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestWithoutDiacritics(t *testing.T) {
	query := "&#233;"
	expect := "b"
	content := fmt.Sprintf(
		"<table><tr><td>%s</td><td>%s</td></tr></table>",
		expect,
		query,
	)

	result, err := GetAdjacent("e", strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestFavorDiacritics(t *testing.T) {
	queryDiacritic := "&#233;"
	query := "é"
	expect := "b"

	content := fmt.Sprintf(
		"<table><tr><td>e</td><td>c</td></tr><tr><td>%s</td><td>%s</td></tr></table>",
		expect,
		queryDiacritic,
	)

	result, err := GetAdjacent(query, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestRemoveDiacritics(t *testing.T) {
	char := "é"
	result := RemoveDiacritics(char)

	if result != "e" {
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	query := "tést"
	expect := "b"

	content := fmt.Sprintf(
		"<table><tr><td>e</td><td>c</td></tr><tr><td>%s</td><td>word1 %s word2</td></tr></table>",
		expect,
		query,
	)

	result, err := GetAdjacent(query, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestContainsNoDiacritics(t *testing.T) {
	query := "tést"
	queryWithoutDiacritics := "test"
	expect := "b"

	content := fmt.Sprintf(
		"<table><tr><td>e</td><td>c</td></tr><tr><td>%s</td><td>word1 %s word2</td></tr></table>",
		expect,
		query,
	)

	result, err := GetAdjacent(queryWithoutDiacritics, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}

func TestTripleRow(t *testing.T) {
	query := "tést"
	queryWithoutDiacritics := "test"
	expect := "b"

	// Only check for the first and last columns. Ignore the middle ones.
	content := fmt.Sprintf(
		("<table>" +
			"<tr><td>e</td><td>%s</td><td>c</td></tr>" +
			"<tr><td>%s</td><td>ignore</td><td>%s</td></tr>" +
			"</table>"),
		query,
		expect,
		query,
	)

	result, err := GetAdjacent(queryWithoutDiacritics, strings.NewReader(content))

	if err != nil || expect != result {
		t.Fail()
	}
}
