package syndication

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type FeedEntry struct {
	Title     string `xml:"title"`
	Subtitle  string `xml:"subtitle"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Author    struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Link    Link   `xml:"link"`
	Content string `xml:"content"`
}

type Feed struct {
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Links    []Link      `xml:"link"`
	Entries  []FeedEntry `xml:"entry"`
}

func ExtractFeedDetails(url string) (*Feed, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching url: %s", resp.Status)
	}

	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	feedMimeTypes := []string{"rss+xml", "atom+xml", "application/xml", "text/xml"}

	var sources []string
	if !slices.Contains(feedMimeTypes, contentType) {
		srcs, err := discoverFeedUrls(url)
		if err != nil {
			return nil, err
		}
		sources = append(sources, srcs...)
	} else {
		sources = append(sources, url)
	}

	if len(sources) == 0 {
		return nil, errors.New("No feed found")
	}

	resp, err = http.Get(sources[0])
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	f := &Feed{}
	xml.NewDecoder(resp.Body).Decode(f)
	return f, nil
}
