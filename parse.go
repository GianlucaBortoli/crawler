package crawler

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func FindLinks(URL string, page []byte) ([]string, error) {
	if len(page) == 0 {
		return []string{}, nil
	}

	reader := bytes.NewReader(page)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to parse page: %v", err)
	}

	var links []string
	// Find all a tags within body
	doc.Find("body a").Each(func(i int, s *goquery.Selection) {
		linkTag := s
		link, _ := linkTag.Attr("href")
		link, err = getAbsURL(URL, link)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		links = append(links, link)
	})
	return links, nil
}

func getAbsURL(from, found string) (string, error) {
	foundURL, err := url.Parse(found)
	if err != nil {
		return "", fmt.Errorf("unable to parse %s: %v", found, err)
	}
	if foundURL.IsAbs() {
		return found, nil
	}
	// Need to build partial URLs starting from the parent page URL
	fromURL, err := url.Parse(from)
	if err != nil {
		return "", fmt.Errorf("unable to parse %s: %v", found, err)
	}

	if foundURL.Scheme == "" {
		foundURL.Scheme = fromURL.Scheme
	}
	if foundURL.Host == "" {
		foundURL.Host = fromURL.Host
	}
	return foundURL.String(), nil
}
