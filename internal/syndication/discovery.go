package syndication

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func extractFeedLink(n *html.Node) *string {
	if n.Type == html.ElementNode && n.Data == "link" {
		var feedType, feedURL string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "type":
				feedType = attr.Val
			case "href":
				feedURL = attr.Val
			}
		}
		if feedType == "application/rss+xml" || feedType == "application/atom+xml" {
			return &feedURL
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := extractFeedLink(c); res != nil {
			return res
		}
	}
	return nil
}

func discoverFeedURL(url string) (*string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return extractFeedLink(doc), nil
}

func isFeedURL(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	contentTypeRaw := resp.Header["Content-Type"]
	if len(contentTypeRaw) == 0 {
		return false, errors.New("Missing content type")
	}

	contentType := strings.Split(contentTypeRaw[0], ";")[0]
	feedMimeTypes := map[string]bool{
		"text/xml":             true,
		"application/xml":      true,
		"application/rss+xml":  true,
		"application/atom+xml": true,
	}

	if contentType == "text/html" {
		return false, nil
	}
	if feedMimeTypes[contentType] {
		return true, nil
	}
	return false, fmt.Errorf("Invalid content type: %s", contentType)
}

func detectFeedType(body []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(body))

	for {
		token, err := decoder.Token()
		if err != nil {
			return "", err
		}
		if start, ok := token.(xml.StartElement); ok {
			return start.Name.Local, nil
		}
	}
}
