package syndication

import "strings"

type AtomFeedEntry struct {
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

func (afe AtomFeedEntry) toFeedEntry() *FeedEntry {
	return &FeedEntry{
		Title:       strings.TrimSpace(afe.Title),
		Description: strings.TrimSpace(afe.Subtitle),
		Published:   afe.Published,
		Updated:     afe.Updated,
		Author:      strings.TrimSpace(afe.Author.Name),
		Link:        afe.Link.Href,
		Content:     strings.TrimSpace(afe.Content),
	}
}

type AtomFeed struct {
	Title   string          `xml:"title"`
	Links   []Link          `xml:"link"`
	Updated string          `xml:"updated"`
	Entries []AtomFeedEntry `xml:"entry"`
}

func (af AtomFeed) toFeed() *Feed {

	var feedURL, siteURL string
	for _, link := range af.Links {
		switch link.Rel {
		case "self":
			feedURL = link.Href
		default:
			siteURL = link.Href
		}
	}
	var entries []FeedEntry
	for _, entry := range af.Entries {
		entries = append(entries, *entry.toFeedEntry())
	}
	return &Feed{
		Title:   af.Title,
		FeedURL: feedURL,
		SiteURL: siteURL,
		Entries: entries,
		Type:    Atom,
	}
}
