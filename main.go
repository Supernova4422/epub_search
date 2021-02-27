package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var port int
	var folder string

	// flags declaration using flag package
	flag.IntVar(&port, "p", 0, "Port this server will run on.")
	flag.StringVar(
		&folder,
		"f",
		"",
		"Path to folder that contains folders that contain source files.",
	)

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		subfolder := r.URL.Query().Get("f")
		query := r.URL.Query().Get("q")
		queryFolder := filepath.Join(folder, subfolder)
		files, err := ioutil.ReadDir(queryFolder)
		if err != nil {
			return
		}

		var bestResult goquery.Selection
		var bestRank int
		found := false

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			fullPath := filepath.Join(queryFolder, file.Name())
			contents, err := os.Open(fullPath)
			defer contents.Close()

			if err != nil {
				return
			}

			rank, result, err := GetAdjacent(query, contents)
			if err == nil && (found == false || rank < bestRank) {
				found = true
				bestRank = rank
				bestResult = *result
			}
		}

		if found {
			fmt.Fprintf(w, "<p id=\"result\">"+bestResult.Text()+"</p>")
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Fatal()
}
