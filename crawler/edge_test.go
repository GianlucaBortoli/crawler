package crawler

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errWriter struct {
	io.Writer
}

func (e errWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("some error")
}

func TestPrintEdge(t *testing.T) {
	e := edge{from: "a", to: "b"}
	var b strings.Builder
	out := io.Writer(&b)

	printEdge(e, out)
	assert.Equal(t, "a, b\n", b.String())
}

func TestPrintEdge_noFrom(t *testing.T) {
	e := edge{to: "b"}
	var b strings.Builder
	out := io.Writer(&b)

	printEdge(e, out)
	assert.Empty(t, b.String())
}

func TestPrintEdge_noTo(t *testing.T) {
	e := edge{from: "a"}
	var b strings.Builder
	out := io.Writer(&b)

	printEdge(e, out)
	assert.Empty(t, b.String())
}

func TestPrintEdge_writeError(t *testing.T) {
	e := edge{from: "a", to: "b"}
	var b errWriter

	printEdge(e, b)
}

func TestPrintSiteMap(t *testing.T) {
	edgesChan := make(chan edge, 1000)
	edges := []edge{{from: "a", to: "b"}, {from: "c", to: "d"}}
	for _, e := range edges {
		edgesChan <- e
	}

	var b strings.Builder
	out := io.Writer(&b)

	PrintSiteMap(edgesChan, out)
	assert.Equal(t, "a, b\nc, d\n", b.String())
}
