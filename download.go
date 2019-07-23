package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Download(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("unable to download page from %s: %v", URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download of %s returned %d status code", URL, resp.StatusCode)
	}
	// RFC 2616 and RFC 7230 define HTTP headers as case insensitive.
	// Transform it in lowercase to make comparison work in all the possible scenarios. 
	contentType := resp.Header.Get("Content-Type")
	if strings.ToLower(contentType) != "text/html; charset=utf-8" {
		return nil, fmt.Errorf("wrong content type %s", contentType)
	}

	return ioutil.ReadAll(resp.Body)
}
