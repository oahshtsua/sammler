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

func (rfe RSSFeedEntry) toFeedEntry() (*FeedEntry, error) {
	published, err := parseRSSDate(rfe.Published)
	if err != nil {
		return nil, err
	}

	return &FeedEntry{
		Title:     strings.TrimSpace(rfe.Title),
		Published: published,
		Link:      strings.TrimSpace(rfe.Link),
		Content:   strings.TrimSpace(rfe.Description),
	}, nil

}

type RSSFeed struct {
	Channel struct {
		Title         string         `xml:"title"`
		Description   string         `xml:"description"`
		LastBuildDate string         `xml:"lastBuildDate"`
		Link          []string       `xml:"link"`
		Items         []RSSFeedEntry `xml:"item"`
		// TODO: find a way to parse both <atom:link> and <link>
		AtomLink string `xml:"atom:link"`
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
		fe, err := entry.toFeedEntry()
		if err != nil {
			continue
		}
		entries = append(entries, *fe)
	}
	return &Feed{
		Title:   rf.Channel.Title,
		FeedURL: rf.Channel.AtomLink,
		SiteURL: rf.Channel.Link[0],
		Entries: entries,
		Type:    RSS,
	}
}
