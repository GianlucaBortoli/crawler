package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	handler := func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	got, err := Download(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, []byte("OK"), got)
}

func TestDownload2(t *testing.T) {
	handler := func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	got, err := Download(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, []byte("OK"), got)
}

func TestDownload_wrongStatusCode(t *testing.T) {
	handler := func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	got, err := Download(ts.URL)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestDownload_wrongContentType(t *testing.T) {
	handler := func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "asd")
		rw.WriteHeader(http.StatusOK)
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	got, err := Download(ts.URL)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestDownload_wrongURL(t *testing.T) {
	got, err := Download("asd")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestIsContentTypeOK(t *testing.T) {
	testCases := []struct {
		ct    string
		expOK bool
	}{
		{
			"text/html",
			true,
		},
		{
			"text/html;",
			true,
		},
		{
			"text/html; charset=utf-8",
			true,
		},
		{
			"text/html; asd",
			true,
		},
		{
			"text/plain",
			false,
		},
	}

	for _, tt := range testCases {
		ok := isContentTypeOK(tt.ct)
		assert.Equal(t, tt.expOK, ok)
	}
}
