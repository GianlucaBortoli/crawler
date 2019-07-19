package crawler

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func FindLinks(page []byte) ([]string, error) {
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
		links = append(links, link)
	})
	return links, nil
}
