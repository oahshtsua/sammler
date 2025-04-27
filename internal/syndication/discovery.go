package syndication

import (
	"net/http"

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
