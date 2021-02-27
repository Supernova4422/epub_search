package main

import (
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func cleanInput(str string) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, ",", " ")
	str = strings.ReplaceAll(str, ".", " ")
	return str
}

// GetAdjacent finds tables, then returns cells adjacent to those matching query.
//
// html is expected to be HTML or XHTML conformant.
// Query is a string.
func GetAdjacent(query string, html io.Reader) (string, error) {
	query = strings.ToLower(query)

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", err
	}

	rows := doc.Find("td")
	rows = rows.FilterFunction(
		func(_ int, col *goquery.Selection) bool {
			return col.Parent().Children().Length() >= 2
		},
	)

	match := foo(rows, func(rhs string) bool {
		rhs = cleanInput(rhs)
		return query == rhs
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	queryWithoutDiacritics := RemoveDiacritics(query)

	match = foo(rows, func(rhs string) bool {
		rhs = cleanInput(rhs)
		return queryWithoutDiacritics == RemoveDiacritics(rhs)
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	match = foo(rows, func(rhs string) bool {
		rhs = cleanInput(rhs)

		for _, word := range strings.Fields(rhs) {
			if word == query {
				return true
			}
		}
		return false
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	match = foo(rows, func(rhs string) bool {
		rhs = cleanInput(rhs)
		for _, word := range strings.Fields(RemoveDiacritics(rhs)) {
			if word == query {
				return true
			}
		}
		return false
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	return "", errors.New("no result found")
}

func foo(rows *goquery.Selection, match func(string) bool) *goquery.Selection {
	return rows.FilterFunction(
		func(_ int, col *goquery.Selection) bool {

			otherCol := col
			if col.Prev().Nodes == nil {
				otherCol = otherCol.Siblings().Last()
			} else if col.Next().Nodes == nil {
				otherCol = otherCol.Siblings().First()
			} else {
				return false
			}

			return match(otherCol.Text())
		},
	)
}

// RemoveDiacritics will remove diacritics from a string.
func RemoveDiacritics(input string) string {
	result, _, _ := transform.String(
		transform.Chain(
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
			norm.NFC,
		),
		input,
	)
	return result
}
