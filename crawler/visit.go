package crawler

import (
	"fmt"
)

func Visit(from string) []edge {
	var edges []edge
	visited := make(map[string]bool)
	queue := NewQueue()

	queue.Push(from)

	for !queue.IsEmpty() {
		l := queue.Pop()
		if visited[l] {
			// Avoid infinite loops in the pages' graph
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
