package syndication

import (
	"net/http"

	"golang.org/x/net/html"
)

func extractFeedLinks(n *html.Node, feeds *[]string) {
	if n.Type == html.ElementNode && n.Data == "link" {
		var feedType, feedUrl string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "type":
				feedType = attr.Val
			case "href":
				feedUrl = attr.Val
			}
		}
		if feedType == "application/rss+xml" || feedType == "application/atom+xml" {
			*feeds = append(*feeds, feedUrl)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractFeedLinks(c, feeds)
	}
}

func discoverFeedUrls(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var feeds []string
	extractFeedLinks(doc, &feeds)

	return feeds, nil
}
