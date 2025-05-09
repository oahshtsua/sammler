package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oahshtsua/sammler/internal/data"
	"github.com/oahshtsua/sammler/internal/syndication"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	version, err := app.queries.GetSQLiteVersion(context.Background())
	if err != nil {
		app.serverError(w, err)
	}
	w.Write([]byte(version))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	unreadEntries, err := app.queries.GetUnreadEntries(context.Background())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "home.html", unreadEntries)
}

func (app *application) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.queries.GetFeeds(context.Background())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "feeds.html", feeds)
}

func (app *application) createFeed(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	url := r.PostForm.Get("feedUrl")
	feedDetails, err := syndication.ExtractFeedDetails(url)
	if err != nil {
		switch {
		case errors.Is(err, syndication.ErrFeedNotFound):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	feedURL, siteURL := extractFeedAndSiteURLs(feedDetails)
	now := time.Now().UTC().Format(time.RFC3339)

	// TODO: wrap feed and entries creation in a transaction

	feed, err := app.queries.CreateFeed(context.Background(), data.CreateFeedParams{
		Title:     feedDetails.Title,
		FeedUrl:   feedURL,
		SiteUrl:   siteURL,
		UpdatedAt: now, // will check and save articles immediately after feed is added
		CheckedAt: now,
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "UNIQUE constraint failed"):
			app.clientError(w, http.StatusConflict)
		default:
			app.serverError(w, err)
		}
		return
	}

	// TODO: figure out a way to insert in bulk with sqlc
	for _, entry := range feedDetails.Entries {
		err := app.queries.CreateEntry(context.Background(), data.CreateEntryParams{
			FeedID:      feed.ID,
			Title:       entry.Title,
			Author:      sql.NullString{String: entry.Author.Name, Valid: entry.Author.Name != ""},
			Content:     "", // TODO: store the description later maybe
			ExternalUrl: entry.Link.Href,
			PublishedAt: entry.Published,
			CreatedAt:   now,
		})
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%d/", feed.ID))
	w.WriteHeader(http.StatusCreated)
}

func (app *application) getFeed(w http.ResponseWriter, r *http.Request) {
	feedID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedID < 1 {
		app.notFound(w)
		return
	}

	feed, err := app.queries.GetFeed(context.Background(), feedID)
	if err != nil {
		switch {
		case err.Error() == "sql: no rows in result set":
			app.notFound(w)
		default:
			app.serverError(w, err)

		}
		return
	}

	entries, err := app.queries.GetFeedEntries(context.Background(), feedID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "feed.html", map[string]any{
		"feed":    feed,
		"entries": entries,
	})
}

func (app *application) deleteFeed(w http.ResponseWriter, r *http.Request) {
	feedID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedID < 1 {
		app.notFound(w)
		return
	}

	err = app.queries.DeleteFeed(context.Background(), feedID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Add("HX-Redirect", "/feeds/")
	w.WriteHeader(http.StatusOK)
}

func (app *application) markFeedRead(w http.ResponseWriter, r *http.Request) {
	feedID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedID < 1 {
		app.notFound(w)
		return
	}

	err = app.queries.MarkFeedRead(context.Background(), feedID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%d/", feedID))
	w.WriteHeader(http.StatusOK)
}

func (app *application) getEntry(w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || entryID < 1 {
		app.notFound(w)
		return
	}

	entry, err := app.queries.GetEntry(context.Background(), entryID)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, http.StatusOK, "entry.html", entry)
}

func (app *application) markEntriesRead(w http.ResponseWriter, r *http.Request) {
	err := app.queries.MarkEntriesRead(context.Background())
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/"))
	w.WriteHeader(http.StatusOK)
}

func (app *application) markEntryRead(w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || entryID < 1 {
		app.notFound(w)
		return
	}

	err = app.queries.MarkEntryRead(context.Background(), entryID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Add("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
