package syndication

import (
	"encoding/xml"
	"net/http"
)

func GetNewEntries(feedURL string, ft FeedType, cutoff string) ([]FeedEntry, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	var newEntries []FeedEntry
	switch ft {
	case RSS:
		updatedFeed := RSSFeed{}
		err := decoder.Decode(&updatedFeed)
		if err != nil {
			return nil, err
		}

		for _, entry := range updatedFeed.Channel.Items {
			fe, err := entry.toFeedEntry()
			if err != nil {
				continue
			}
			if fe.Published < cutoff {
				break
			}
			newEntries = append(newEntries, *fe)
		}
	case Atom:
		updatedFeed := AtomFeed{}
		err := decoder.Decode(&updatedFeed)
		if err != nil {
			return nil, err
		}
		for _, entry := range updatedFeed.Entries {
			fe := entry.toFeedEntry()
			if fe.Published < cutoff {
				break
			}
			newEntries = append(newEntries, *fe)
		}
	default:
		return nil, ErrFeedNotSupported

	}
	return newEntries, nil
}
