package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisit(t *testing.T) {
	links := Visit("https://monzo.com")
	assert.NotEqual(t, 0, len(links))
}

func TestIsSameSubDomain(t *testing.T) {
	testCases := []struct {
		a             string
		b             string
		sameSubDomain bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"98:a//",
			"",
			false,
		},
		{
			"",
			"98:a//",
			false,
		},
		{
			"98:a//",
			"98:a//",
			false,
		},
		{
			"http://asd.com",
			"http://asd.com",
			true,
		},
		{
			"http://asd.com",
			"https://asd.com",
			true,
		},
		{
			"http://asd.com/foo/bar",
			"http://asd.com",
			true,
		},
		{
			"http://asd.com/foo&bar",
			"http://asd.com",
			true,
		},
		{
			"http://www.asd.com",
			"http://asd.com",
			true,
		},
		{
			"http://www.asd.com",
			"https://asd.com",
			true,
		},
		{
			"http://foo.asd.com",
			"http://asd.com",
			false,
		},
		{
			"http://asd.com:443",
			"http://asd.com:80",
			true,
		},
	}

	for _, tt := range testCases {
		res := isSameSubDomain(tt.a, tt.b)
		assert.Equal(t, tt.sameSubDomain, res)
	}
}
