# crawler
Simple web crawler in Go

## Description
This is a simple website crawler. You can pass it a URL and it will find all the links it can
reach as long as they're in the same sub-domain.
A textual sitemap is printed in the end to show where the URL visit reached.

## Build
The project uses go modules and requires go `>= 1.12`. To build the binary run
```bash
make build
```

The executable will be placed in `./bin` and has the race detector enabled.

## Unit tests
Unit tests can be run (with the race detector enabled) typing 
```bash
make test
```

The overall code coverage is also outputted.

## Run
The only required parameter is the URL where it needs to start crawling from.
For example:
```bash
./bin/crawler -url http://example.com
```

Other parameters can be configured and has default values. Run
```bash
./bin/crawler [-h]
```
to see all the possibilities.

## Design decisions & limitations
During the development of the web crawler, some design decisions and trade-offs have been made.

Here is a list of the main ones.

* Download and parse:
    * download only via HTTP GET method
    * no retry mechanism
    * consider only links whose advertised `Content-Type` HTTP header is of the `text/html` family
    * follow only links which are inside the `href` attribute of `<a>` tags (eg.
    `<a href="www.example.com">Link</a>`)
    * redirects, timeouts etc follows the default of the HTTP client of the golang standard library
    (see https://golang.org/pkg/net/http/#Client)
    * not respecting `robots.txt`
* Visit:
    * keep track of already visited URLs to avoid loops
    * configurable max depth for the visit, since we cannot know in advance the graph size
    * implemented as a breadth-first traversal so that multiple "workers" can concurrently download
    pages and find links inside them
* Output:
    * visit information logged to standard error
    * site map printed to standard output with a `<from>, <to>` format

## Improvements
Given the design decisions of the previous section, it's possible to think of possible improvements
of the current crawler implementation.

* The current implementation uses a map `URL -> visited` to keep track of URLs that have already
been visited. This is a simple and effective solution for most of the cases, but there are some where
the crawler can visit the same page multiple times (eg. when two different URLs point to the same page).
To avoid this problem, we need to change the way it detects possible loops. A possible improvement
could be to use both the URL and the page content, computing some degree of similarity with the new
URLs that are being visited. In this way the crawler should decide when two pages are "similar enough".
* Split the "download page" and "parse page to find links" logic to different workers. Right now the
same worker handles them sequentially, but building a more sophisticated data pipeline can help to
spread the work more fairly and make the crawler more efficient. Normally parsing the HTML of a page
is faster than downloading one, so it would make sense to have more workers that download pages rather
than workers that parse them.
* Every worker could reuse the same `http.Transport` (https://godoc.org/net/http#Transport) for the
HTTP client. This allows to reuse HTTP connections and TCP sockets that have already been opened before,
limiting also the number of TLS handshakes (which are quite expensive in terms of time). This can be
of help especially when most of the visited URLs are within the same sub-domain.
* Respect `robots.txt` advertised by web servers to respect site owners' wills.
