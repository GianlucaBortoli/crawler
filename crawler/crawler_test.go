package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c, e := New("asd", 100000)
	assert.NotNil(t, c)
	assert.IsType(t, &Crawler{}, c)
	assert.NotNil(t, e)
}

func TestCrawler_Start(t *testing.T) {
	c, _ := New("https://google.com", 5)
	c.Start()
	c.Wait()
}

func TestCrawler_StartManyTimes(t *testing.T) {
	c, _ := New("https://google.com", 5)
	c.Start()
	c.Start()
	c.Start()
	c.Wait()
}

func TestCrawler_StartWaitStop(t *testing.T) {
	c, _ := New("https://google.com", 5)

	c.Start()
	c.stop()
}

func TestVisitURL(t *testing.T) {
	testCases := []struct {
		URL          string
		shouldErr    bool
		shouldResNil bool
	}{
		{
			"https://google.com",
			false,
			false,
		},
		{
			"https://cnn.com",
			false,
			false,
		},
		{
			"https://asd.asd",
			true,
			true,
		},
		{
			"asd",
			true,
			true,
		},
		{
			"",
			true,
			true,
		},
	}

	for _, tt := range testCases {
		res, err := visitURL(tt.URL)
		assert.Equal(t, tt.shouldErr, err != nil)
		assert.Equal(t, tt.shouldResNil, res == nil)
	}
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

func TestIsAlreadyVisited(t *testing.T) {
	c, _ := New("asd", 5)

	got := c.isAlreadyVisited("asd")
	assert.False(t, got)

	c.visited.Store("asd", true)
	got = c.isAlreadyVisited("asd")
	assert.True(t, got)
}

func TestSetVisited(t *testing.T) {
	c, _ := New("asd", 5)

	_, visited := c.visited.Load("asd")
	assert.False(t, visited)

	c.setVisited("asd")
	_, visited = c.visited.Load("asd")
	assert.True(t, visited)
}
