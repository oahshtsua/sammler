package syndication

type FeedType string

const (
	RSS  FeedType = "rss"
	Atom FeedType = "atom"
)

type FeedConvertible interface {
	toFeed() *Feed
}
type FeedEntry struct {
	Title       string
	Description string
	Published   string
	Updated     string
	Author      string
	Link        string
	Content     string
}

type Feed struct {
	Type    FeedType
	Title   string
	FeedURL string
	SiteURL string
	Entries []FeedEntry
}

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr"`
}
