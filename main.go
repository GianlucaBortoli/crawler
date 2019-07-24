package main

import (
	"fmt"
	"os"

	"github.com/cog-qlik/crawler/crawler"
)

func help() {
	fmt.Println(`Usage:
crawler <url>

Example:
crawler https://monzo.com`)
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Wrong number of arguments")
		help()
	}

	link := args[1]
	edges := crawler.Visit(link)
	for _, e := range edges {
		e.Print()
	}
}
