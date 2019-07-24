package crawler

import (
	"fmt"
)

type edge struct {
	from string
	to   string
}

func (e edge) Print() {
	fmt.Printf("%s -> %s\n", e.from, e.to)
}

func Visit(from string) []edge {
	var edges []edge
	visited := make(map[string]bool)
	queue := NewQueue()

	queue.Push(from)

	for !queue.IsEmpty() {
		l := queue.Pop()
		if visited[l] {
			continue
		}

		body, downloadErr := Download(l)
		if downloadErr != nil {
			fmt.Println("ERROR:", downloadErr)
		}
		to, findErr := FindLinks(l, body)
		if findErr != nil {
			fmt.Println("ERROR:", findErr)
		}

		// Add edge from starting URL to all extracted links
		for _, l := range to {
			edges = append(edges, edge{from: from, to: l})
			queue.Push(l)
			visited[l] = true
		}
	}
	return edges
}
