package syndication

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

var ErrFeedNotFound = errors.New("No feed found for given URL")
var ErrFeedNotSupported = errors.New("Unsupported feed type")

func resolveFeedURL(url string) (string, error) {
	isFeed, err := isFeedURL(url)
	if err != nil {
		return "", err
	}
	if isFeed {
		return url, nil
	}

	source, err := discoverFeedURL(url)
	if err != nil {
		return "", err
	}
	if source == nil {
		return "", ErrFeedNotFound
	}
	return *source, nil
}

func parseFeed(data []byte, feedURL string) (*Feed, error) {
	feedType, err := detectFeedType(data)
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	switch feedType {
	case "rss":
		f := RSSFeed{}
		err := decoder.Decode(&f)
		if err != nil {
			return nil, err
		}
		// HACK: manually assign the feed link to the Feed item
		f.Channel.AtomLink = feedURL
		return f.toFeed(), nil
	case "feed":
		f := AtomFeed{}
		err := decoder.Decode(&f)
		if err != nil {
			return nil, err
		}
		return f.toFeed(), nil
	default:
		return nil, ErrFeedNotSupported
	}
}

func ExtractFeedDetails(url string) (*Feed, error) {

	source, err := resolveFeedURL(url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseFeed(body, source)

}
