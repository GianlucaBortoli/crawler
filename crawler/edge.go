package crawler

import (
	"fmt"
	"io"
)

type edge struct {
	from string
	to   string
}

func PrintSiteMap(edgesChan <-chan edge, w io.Writer) {
	for {
		select {
		case e := <-edgesChan:
			printEdge(e, w)
		default:
			return
		}
	}
}

func printEdge(e edge, w io.Writer) {
	if e.from == "" || e.to == "" {
		fmt.Println("ERROR: malformed edge")
		return
	}

	s := fmt.Sprintf("%s, %s\n", e.from, e.to)
	if _, err := io.WriteString(w, s); err != nil {
		fmt.Println("ERROR: unable to write edge")
	}
}
