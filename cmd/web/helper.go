package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/oahshtsua/sammler/internal/data"
	"github.com/oahshtsua/sammler/internal/syndication"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.logger.Error(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func extractFeedAndSiteURLs(details *syndication.Feed) (string, string) {
	var feedURL, siteURL string
	for _, link := range details.Links {
		switch link.Rel {
		case "self":
			feedURL = link.Href
		default:
			siteURL = link.Href
		}
	}
	return feedURL, siteURL
}

func (app *application) refreshFeeds() error {
	feeds, err := app.queries.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		newEntries, err := syndication.GetNewEntries(feed)
		if err != nil {
			app.logger.Error(fmt.Sprintf("Error refreshing feed: %s", feed.Title))
			return err
		}

		now := time.Now().UTC().Format(time.RFC3339)
		for _, entry := range newEntries {
			err := app.queries.CreateEntry(context.Background(), data.CreateEntryParams{
				FeedID:      feed.ID,
				Title:       entry.Title,
				Author:      sql.NullString{String: entry.Author.Name, Valid: entry.Author.Name != ""},
				Content:     "",
				ExternalUrl: entry.Link.Href,
				PublishedAt: entry.Published,
				CreatedAt:   now,
			})

			if err != nil {
				return err
			}

		}
		err = app.queries.UpdateFeedCheckedAt(context.Background(), data.UpdateFeedCheckedAtParams{
			ID:        feed.ID,
			CheckedAt: now,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
