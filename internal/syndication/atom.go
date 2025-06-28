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
		entries = append(entries,
			FeedEntry{
				Title:       strings.TrimSpace(entry.Title),
				Description: strings.TrimSpace(entry.Subtitle),
				Published:   entry.Published,
				Updated:     entry.Updated,
				Author:      strings.TrimSpace(entry.Author.Name),
				Link:        entry.Link.Href,
				Content:     strings.TrimSpace(entry.Content),
			})
	}
	return &Feed{
		Title:   af.Title,
		FeedURL: feedURL,
		SiteURL: siteURL,
		Entries: entries,
		Type:    Atom,
	}
}
