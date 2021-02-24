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

/*

func FindMatch(query string, html io.Reader, removeDiacritics bool) (string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	res, err := GetAdjacent(query, doc, false)
	if err == nil {
		return res, nil
	}
}
*/

// GetAdjacent finds tables, then returns cells adjacent to those matching query.
//
// html is expected to be HTML or XHTML conformant.
// Query is a string.
func GetAdjacent(query string, html io.Reader, removeDiacritics bool) (string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", err
	}

	if removeDiacritics {
		query = RemoveDiacritics(query)
	}

	rows := doc.Find("td")
	rows = rows.FilterFunction(
		func(_ int, col *goquery.Selection) bool {
			return col.Parent().Children().Length() == 2
		},
	)

	match := foo(rows, func(rhs string) bool {
		return query == rhs
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	queryWithoutDiacritics := RemoveDiacritics(query)

	match = foo(rows, func(rhs string) bool {
		return queryWithoutDiacritics == RemoveDiacritics(rhs)
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	match = foo(rows, func(rhs string) bool {
		return strings.Contains(rhs, query)
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	match = foo(rows, func(rhs string) bool {
		return strings.Contains(RemoveDiacritics(rhs), queryWithoutDiacritics)
	})
	if len(match.Nodes) != 0 {
		return match.Text(), nil
	}

	return "", errors.New("no result found")
}

func foo(rows *goquery.Selection, match func(string) bool) *goquery.Selection {
	return rows.FilterFunction(
		func(_ int, col *goquery.Selection) bool {
			otherCol := col.Next()
			if len(otherCol.Nodes) == 0 {
				otherCol = col.Prev()
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
