package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var port int
	var folder string

	// flags declaration using flag package
	flag.IntVar(&port, "p", 0, "Port this server will run on.")
	flag.StringVar(&folder, "f", "", "Specify folder that contians source files.")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			fullPath := filepath.Join(folder, file.Name())
			contents, err := os.Open(fullPath)
			defer contents.Close()

			if err != nil {
				return
			}

			result, err := GetAdjacent(query, contents)
			if err == nil {
				fmt.Fprintf(w, "<p id=\"result\">"+result+"</p>")
			}
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Fatal()
}
