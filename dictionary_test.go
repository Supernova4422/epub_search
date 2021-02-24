package main

import (
	"fmt"
	"strings"
	"testing"
)

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
