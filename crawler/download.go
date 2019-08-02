package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Download downloads a web page for the given URL using the GET method.
// This function returns an error if either:
//   * the returned status code is "200 OK"
//   * the returned Content-Type is of the "text/html" family
func Download(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("unable to download page from %s: %v", URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download of %s returned %d status code", URL, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !isContentTypeOK(contentType) {
		return nil, fmt.Errorf("wrong content type %s", contentType)
	}
	return ioutil.ReadAll(resp.Body)
}

func isContentTypeOK(s string) bool {
	// RFC 2616 and RFC 7230 define HTTP headers as case insensitive.
	// Transform it in lowercase to make comparison work in all the possible scenarios.
	return strings.Contains(strings.ToLower(s), "text/html")
}
