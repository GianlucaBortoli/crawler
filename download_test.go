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
