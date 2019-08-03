package crawler

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// Crawler implements the web crawler with many workers that can
// visit websites
type Crawler struct {
	workers  int
	maxDepth int
	log      *log.Logger

	visited   sync.Map    // stores URLs which are already visited (URL -> bool)
	URLChan   chan string // input
	edgesChan chan edge   // output
	quitChan  chan struct{}
	wg        sync.WaitGroup
	iter      int32
	startOnce sync.Once
	stopOnce  sync.Once
}

// New creates a crawler that starts visiting websites the given URL with
// a given amount of workers.
// Returns the crawler and a channel where edges are sent as the web pages
// are visited
func New(URL string, workers, maxDepth int) (*Crawler, <-chan edge) {
	edgesChan := make(chan edge, 10000)
	URLChan := make(chan string, 10000)
	// Enqueue starting URL
	URLChan <- URL

	return &Crawler{
		workers:   workers,
		maxDepth:  maxDepth,
		log:       log.New(os.Stderr, "", log.LstdFlags),
		URLChan:   URLChan,
		edgesChan: edgesChan,
		quitChan:  make(chan struct{}, workers),
	}, edgesChan
}

// Start starts the crawling procedure
func (c *Crawler) Start() {
	// Ensure workers can be started only once to avoid leaking goroutines
	c.startOnce.Do(func() {
		for i := 0; i < c.workers; i++ {
			c.wg.Add(1)
			go c.start()
		}
	})
}

// stop gracefully stops every worker in the crawler.
// This is used only in unit-test only for now.
func (c *Crawler) stop() { //nolint:unused
	// Ensure workers are stopped only once to avoid sending too many signals in the quitChan
	c.stopOnce.Do(func() {
		for i := 0; i < c.workers; i++ {
			c.quitChan <- struct{}{}
		}
	})
}

// Wait waits until the crawling procedure ends. This can be useful for printing
// the site map when the visit is completed
func (c *Crawler) Wait() {
	c.wg.Wait()
}

// start starts the URL visit process in a breadth-first manner.
// The URLChan acts like a queue for pushing/popping nodes during the visit.
func (c *Crawler) start() {
	defer c.wg.Done()

	for {
		select {
		case from := <-c.URLChan:
			if c.isMaxDepth() {
				c.log.Printf("[WARN] Max iterations %d reached\n", c.iter)
				return
			}
			// Skip URLs that have already been visited before. We want to avoid possible
			// infinite loops in the graph
			if c.isAlreadyVisited(from) {
				c.log.Printf("[INFO] link %s already visited\n", from)
				continue
			}

			to, err := visitURL(from)
			if err != nil {
				c.log.Println("[ERROR]", err)
			}
			c.setVisited(from)
			c.enqueueChildren(from, to)

			atomic.AddInt32(&c.iter, 1)
		case <-c.quitChan:
			return
		default:
			return
		}
	}
}

func (c *Crawler) isMaxDepth() bool {
	return atomic.LoadInt32(&c.iter) >= int32(c.maxDepth)
}

func (c *Crawler) setVisited(URL string) {
	c.visited.Store(URL, true)
}

func (c *Crawler) isAlreadyVisited(URL string) bool {
	_, ok := c.visited.Load(URL)
	return ok
}

func (c *Crawler) enqueueChildren(from string, to []string) {
	for _, t := range to {
		// Don't follow links in a different sub-domain
		if !isSameSubDomain(from, t) {
			c.log.Printf("[INFO] %s not in the subdomain of %s. Skipping\n", t, from)
			continue
		}
		// Enqueue new node
		c.URLChan <- t
		// Send edges as I visit them
		c.edgesChan <- edge{from: from, to: t}
	}
}

// visitURL downloads the page and finds all the links inside
func visitURL(URL string) ([]string, error) {
	body, downloadErr := Download(URL)
	if downloadErr != nil {
		return nil, fmt.Errorf("unable to download %s: %v", URL, downloadErr)
	}
	to, findErr := FindLinks(URL, body)
	if findErr != nil {
		return nil, fmt.Errorf("unable to find children links for %s: %v", URL, findErr)
	}
	return to, nil
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
