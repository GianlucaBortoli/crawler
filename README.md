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

## Design decisions
Some design decisions and trade-offs have been made during the development of this web crawler.
Here is a list of the main ones divided by topics.

* Page download and parse:
    * Download only via HTTP GET method as a browser would do.
    * No retry mechanism as a browser would do.
    * Consider only links whose advertised `Content-Type` HTTP header is of the `text/html` family.
    In this way the crawler reads the response body only when its contend is advertised as an HTML page.
    * Follow only links which are inside the `href` attribute of `<a>` tags (eg.
    `<a href="www.example.com">Link</a>`).
    * Redirects, timeouts etc follows the default of the HTTP client of the golang standard library
    (see https://golang.org/pkg/net/http/#Client).
    * Not respecting `robots.txt`.
* Visit:
    * Keep track of already visited URLs to avoid loops.
    * Configurable max depth for the visit, since we cannot know in advance the graph size.
    * Implemented as a breadth-first traversal. In this way multiple workers can concurrently download
    pages and find links inside them.
* Output:
    * Visit information logged to standard error to see the progress.
    * Site map printed to standard output with a `<from>, <to>` format. This allows to pipe the
    site map (and not the logs) where we want, possibly becoming the input of another program.

## Improvements
Given the decisions listed in the previous section, it's possible to think of possible improvements
that can enhance the overall performance.

* The current implementation uses a `map[string]bool` structure (`URL -> visited`) to keep track of
URLs that have already been visited.
This is a simple and effective solution for most of the cases, but there are some where the crawler
can visit the same page multiple times (eg. when two different URLs point to the same page content).
To avoid this problem, we need to improve the loop detection phase. A possible improvement
could be to use both the URL string and the page content, computing some degree of similarity with the new
pages that are being visited. In this way the crawler has more context to decide when two pages
are "similar enough".
* Split the "download page" and "parse page to find links" logic into different workers. Right now the
same worker handles them sequentially, but building a more sophisticated data pipeline can help to
spread the work more fairly and make the entire process more efficient. Normally, parsing the HTML of a page
is faster than downloading one, so it would make sense to have more workers that download pages rather
than workers that parse them.
* Every worker could reuse the same `http.Transport` (https://godoc.org/net/http#Transport) for the
HTTP client. This allows to reuse HTTP connections and TCP sockets that have already been opened before,
limiting also the number of TLS handshakes (which are quite expensive in terms of time). This can help
especially when a high number of workers is used, limiting the usage of file descriptors, and when many
of the URLs belong to the same sub-domain.
* Respect `robots.txt` advertised by web servers to respect site owners' wills.
