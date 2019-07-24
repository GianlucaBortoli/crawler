package crawler

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintEdge(t *testing.T) {
	e := edge{from: "a", to: "b"}
	var b strings.Builder
	out := io.Writer(&b)

	PrintEdge(e, out)
	assert.Equal(t, "a -> b\n", b.String())
}

func TestPrintEdge_noFrom(t *testing.T) {
	e := edge{to: "b"}
	var b strings.Builder
	out := io.Writer(&b)

	PrintEdge(e, out)
	assert.Empty(t, b.String())
}

func TestPrintEdge_noTo(t *testing.T) {
	e := edge{from: "a"}
	var b strings.Builder
	out := io.Writer(&b)

	PrintEdge(e, out)
	assert.Empty(t, b.String())
}
