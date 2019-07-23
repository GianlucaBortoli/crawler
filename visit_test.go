package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisit(t *testing.T) {
	links := Visit("https://monzo.com")
	assert.NotEqual(t, 0, len(links))
}
