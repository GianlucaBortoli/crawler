package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLinks(t *testing.T) {
	testCases := []struct {
		page     string
		expLinks []string
		expError error
	}{
		{
			`<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`,
			[]string{"foo", "/bar/baz"},
			nil,
		},
		{
			"",
			[]string{},
			nil,
		},
		{
			"asd",
			nil,
			nil,
		},
		{
			`<asd></c>`,
			nil,
			nil,
		},
	}

	for _, tt := range testCases {
		links, err := FindLinks([]byte(tt.page))
		assert.Equal(t, tt.expError, err)
		assert.Equal(t, tt.expLinks, links)
	}
}
