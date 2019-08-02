package main

import (
	"flag"
	"os"

	"github.com/cog-qlik/crawler/crawler"
)

func main() {
	URL := flag.String("url", "", "The URL to start crawling from")
	workers := flag.Int("workers", 5, "The number of workers to visit URLs")
	depth := flag.Int("maxdepth", 10, "The maximum depth of the visit")
	flag.Parse()

	if *URL == "" {
		flag.Usage()
		os.Exit(1)
	}

	cr, edges := crawler.New(*URL, *workers, *depth)
	cr.Start()
	cr.Wait()
	crawler.PrintSiteMap(edges, os.Stdout)
}
