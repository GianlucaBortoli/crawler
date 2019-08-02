package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type Crawler struct {
	workers   int
	visited   sync.Map
	URLChan   chan string // input
	edgesChan chan edge   // output
	wg        sync.WaitGroup
}

func New(URL string, workers int) (*Crawler, <-chan edge) {
	edgesChan := make(chan edge, 1000)
	URLChan := make(chan string, 1000)
	// Enqueue starting URL
	URLChan <- URL

	return &Crawler{
		workers:   workers,
		URLChan:   URLChan,
		edgesChan: edgesChan,
	}, edgesChan
}

func (c *Crawler) Start() {
	for i := 0; i < c.workers; i++ {
		c.wg.Add(1)
		go c.start()
	}
}

func (c *Crawler) Wait() {
	c.wg.Wait()
}

func (c *Crawler) start() {
	for {
		select {
		case from := <-c.URLChan:
			if _, ok := c.visited.Load(from); ok {
				// Avoid infinite loops in the pages' graph
				fmt.Printf("[INFO] link %s already visited\n", from)
				continue
			}
			to := visitURL(from)
			c.enqueueChildren(from, to)
		default:
			c.wg.Done()
			return
		}
	}
}

func (c *Crawler) enqueueChildren(from string, to []string) {
	for _, t := range to {
		if !isSameSubDomain(from, t) {
			// Avoid scraping links not in the sub-domain of the initial URL
			fmt.Printf("[INFO] %s not in the subdomain of %s. Skipping\n", t, from)
			continue
		}

		c.URLChan <- t
		c.visited.Store(t, true)
		c.edgesChan <- edge{from: from, to: t}
	}
}

func visitURL(URL string) []string {
	body, downloadErr := Download(URL)
	if downloadErr != nil {
		fmt.Println("[ERROR]", downloadErr)
	}
	to, findErr := FindLinks(URL, body)
	if findErr != nil {
		fmt.Println("[ERROR]", findErr)
	}
	return to
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
