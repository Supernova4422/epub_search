package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

// GetAdjacent finds tables, then returns cells adjacent to those matching query.
//
// html is expected to be HTML or XHTML conformant.
// Query is a string.
func GetAdjacent(query string, html io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", err
	}

	rows := doc.Find("tr")
	result := ""
	found := errors.New("no result found")
	rows = rows.EachWithBreak(
		func(_ int, row *goquery.Selection) bool {
			cols := row.Children()

			firstText := cols.First().Text()
			secondText := cols.Last().Text()

			if firstText == query {
				fmt.Print(secondText)
				result = secondText
				found = nil
				return true
			}

			if secondText == query {
				result = firstText
				found = nil
				return true
			}

			return true
		},
	)

	return result, found
}
