package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

		bestRank := -1
		matches := make([]string, 0)

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
			if err == nil {
				if bestRank < 0 {
					bestRank = rank
				}

				for i := 0; i < len(result.Nodes); i++ {
					node := result.Eq(i)
					parent := node.Parent()
					output := parent.Text()
					fmt.Print(output)
					output = strings.ReplaceAll(output, "\n", " ")
					if rank < bestRank {
						matches = append(matches, "")
						copy(matches[1:], matches)
						matches[0] = output
					} else {
						matches = append(matches, output)
					}
				}
			}
		}

		result := ""
		for _, match := range matches {
			joined := result
			if joined != "" {
				joined += "\n\n"
			}
			
			joined += match
			if len(joined) > 1200 {
				break
			} else {
				result = joined
			}
		}
		fmt.Fprintf(w, "<p id=\"result\">"+result+"</p>")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Fatal()
}
