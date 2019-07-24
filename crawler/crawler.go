package crawler

import (
	"fmt"
	"sync"
)

type Crawler struct {
	visited   sync.Map
	URLChan   chan string // input
	edgesChan chan edge   // output
	quitChan  chan struct{}
}

func New(URL string) (*Crawler, <-chan edge) {
	edgesChan := make(chan edge)
	URLChan := make(chan string, 1000)
	URLChan <- URL

	return &Crawler{
		URLChan:   URLChan,
		edgesChan: edgesChan,
		quitChan:  make(chan struct{}),
	}, edgesChan
}

func (c *Crawler) Start() {
	go c.start()
}

func (c *Crawler) Stop() {
	c.quitChan <- struct{}{}
	close(c.edgesChan)
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
			to := visitSinglePage(from)
			c.sendEdges(from, to)
		case <-c.quitChan:
			fmt.Println("[INFO] quit signal received")
			break
		default:
			//fmt.Println("[INFO] nothing in channels")
			break
		}
	}
}

func (c *Crawler) sendEdges(from string, to []string) {
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

func visitSinglePage(URL string) []string {
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
