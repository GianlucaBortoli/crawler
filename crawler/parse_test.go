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
		links, err := FindLinks("asd", []byte(tt.page))
		assert.Equal(t, tt.expError, err)
		assert.Equal(t, tt.expLinks, links)
	}
}

func TestGetAbsURL(t *testing.T) {
	testCases := []struct {
		from     string
		found    string
		exp      string
		expError bool
	}{
		{
			"http://asd.com",
			"http://asd.com/foo",
			"http://asd.com/foo",
			false,
		},
		{
			"http://asd.com",
			"http://asd.com/foo/bar",
			"http://asd.com/foo/bar",
			false,
		},
		{
			"http://asd.com",
			"/foo",
			"http://asd.com/foo",
			false,
		},
		{
			"http://asd.com",
			"/foo/bar",
			"http://asd.com/foo/bar",
			false,
		},
		{
			"http://asd.com",
			":asda::a",
			"",
			true,
		},
		{
			":asda::a",
			"/foo",
			"",
			true,
		},
	}

	for _, tt := range testCases {
		absURL, err := getAbsURL(tt.from, tt.found)
		assert.Equal(t, err != nil, tt.expError)
		assert.Equal(t, tt.exp, absURL)
	}
}
