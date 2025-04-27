package syndication

import (
	"encoding/xml"
	"net/http"

	"github.com/oahshtsua/sammler/internal/data"
)

func GetNewEntries(f data.Feed) ([]*FeedEntry, error) {
	resp, err := http.Get(f.FeedUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	updatedFeed := &Feed{}
	err = xml.NewDecoder(resp.Body).Decode(updatedFeed)
	if err != nil {
		return nil, err
	}

	newEntries := []*FeedEntry{}
	for _, entry := range updatedFeed.Entries {
		if entry.Published < f.CheckedAt {
			break
		}
		newEntries = append(newEntries, &entry)
	}

	return newEntries, nil
}
