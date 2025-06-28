package syndication

import (
	"fmt"
	"strings"
	"time"
)

type RSSFeedEntry struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Published   string `xml:"pubDate"`
	Link        string `xml:"link"`
}

type RSSFeed struct {
	Channel struct {
		Title         string         `xml:"title"`
		Description   string         `xml:"description"`
		LastBuildDate string         `xml:"lastBuildDate"`
		Link          string         `xml:"link"`
		AtomLink      Link           `xml:"http://www.w3.org/2005/Atom link"`
		Items         []RSSFeedEntry `xml:"item"`
	} `xml:"channel"`
}

func parseRSSDate(date string) (string, error) {
	dateFormats := []string{
		time.RFC1123,
		time.RFC1123Z,
	}
	for _, format := range dateFormats {
		t, err := time.Parse(format, date)
		if err == nil {
			return t.Format(time.RFC3339), nil
		}
	}
	return "", fmt.Errorf("Unrecognized date format: %s", date)
}

func (rf RSSFeed) toFeed() *Feed {
	var entries []FeedEntry
	for _, entry := range rf.Channel.Items {
		published, err := parseRSSDate(entry.Published)
		if err != nil {
			// log the error
			continue
		}
		entries = append(entries,
			FeedEntry{
				Title:     strings.TrimSpace(entry.Title),
				Published: published,
				Link:      strings.TrimSpace(entry.Link),
				Content:   strings.TrimSpace(entry.Description),
			})
	}
	return &Feed{
		Title:   rf.Channel.Title,
		FeedURL: rf.Channel.AtomLink.Href,
		SiteURL: rf.Channel.Link,
		Entries: entries,
		Type:    RSS,
	}
}
