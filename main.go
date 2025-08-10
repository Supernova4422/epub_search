package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

func main() {
	log.Println("Starting")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

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

	log.Println("Using Port:", port, "with folder:", folder)



	var (
		opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
				Name: "epubSearch_Query",
				Help: "The total number of processed queries",
		})
	)

	tracer := otel.Tracer("epub_search")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opsProcessed.Inc()

		_, span := tracer.Start(ctx, "request")
		defer span.End()

		subfolder := r.URL.Query().Get("f")
		query := r.URL.Query().Get("q")

		log.Println("Received request with query:", query)

		queryFolder := filepath.Join(folder, subfolder)
		files, err := ioutil.ReadDir(queryFolder)
		if err != nil {
			return
		}

		bestRank := -1
		matches := make([]string, 0)

		for _, file := range files {
			_, fileSpan := tracer.Start(ctx, "search" + file.Name())

			if file.IsDir() {
				continue
			}

			fullPath := filepath.Join(queryFolder, file.Name())
			contents, err := os.Open(fullPath)
			defer contents.Close()

			if err != nil {
				fileSpan.End()
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

			fileSpan.End()
		}

		result := ""
		for _, match := range matches {
			joined := result
			if joined != "" {
				joined += "\n"
			}
			
			joined += match
			if len(joined) > 1500 {
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
