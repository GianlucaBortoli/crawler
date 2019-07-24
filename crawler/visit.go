package crawler

import (
	"fmt"
	"net/url"
	"strings"
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
			fmt.Printf("[INFO] link %s already visited\n", l)
			continue
		}

		body, downloadErr := Download(l)
		if downloadErr != nil {
			fmt.Println("[ERROR] ", downloadErr)
		}
		to, findErr := FindLinks(l, body)
		if findErr != nil {
			fmt.Println("[ERROR] ", findErr)
		}

		// Add edge from starting URL to all extracted links
		for _, l := range to {
			if !isSameSubDomain(from, l) {
				// Avoid scraping links not in the sub-domain of the initial URL
				fmt.Printf("[INFO] %s not in the subdomain of %s. Skipping\n", l, from)
				continue
			}
			edges = append(edges, edge{from: from, to: l})
			queue.Push(l)
			visited[l] = true
		}
	}
	return edges
}

func isSameSubDomain(a, b string) bool {
	aParsed, aErr := url.Parse(a)
	if aErr != nil {
		return false
	}
	bParsed, bErr := url.Parse(b)
	if bErr != nil {
		return false
	}
	// Host fields can start with "www.". They don't make any difference in the sub-domain
	// so trim the prefix.
	// The Hostname() function already takes care of stripping the port
	// from the host. I want https://asd:80 and https://asd:443 to be in the same sub-domain
	aDomain := strings.TrimPrefix(aParsed.Hostname(), "www.")
	bDomain := strings.TrimPrefix(bParsed.Hostname(), "www.")
	return aDomain == bDomain
}
