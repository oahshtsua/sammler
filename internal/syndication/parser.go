package syndication

import (
	"encoding/xml"
	"errors"
	"net/http"
)

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type FeedEntry struct {
	Title     string `xml:"title"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Author    struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Link Link `xml:"link"`
}

type Feed struct {
	Title   string      `xml:"title"`
	Links   []Link      `xml:"link"`
	Entries []FeedEntry `xml:"entry"`
}

var ErrFeedNotFound = errors.New("No feed found for given URL")

func ExtractFeedDetails(url string) (*Feed, error) {
	source, err := discoverFeedURL(url)
	if source == nil {
		return nil, ErrFeedNotFound
	}

	resp, err := http.Get(*source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	f := &Feed{}
	xml.NewDecoder(resp.Body).Decode(f)
	return f, nil
}
