package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
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

	now := time.Now().UTC().Format(time.RFC3339)

	feed, err := app.queries.CreateFeed(context.Background(), data.CreateFeedParams{
		Title:     feedDetails.Title,
		Type:      feedDetails.Type,
		FeedUrl:   feedDetails.FeedURL,
		SiteUrl:   feedDetails.SiteURL,
		UpdatedAt: now,
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

	entries := buildCreateEntryParams(feed.ID, now, feedDetails.Entries)
	err = app.queries.CreateMultipleEntry(context.Background(), entries)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%d/", feed.ID))
	w.WriteHeader(http.StatusCreated)
}

func (app *application) getFeed(w http.ResponseWriter, r *http.Request) {
	feedID, err := parseID(r)
	if err != nil {
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
	feedID, err := parseID(r)
	if err != nil {
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
	feedID, err := parseID(r)
	if err != nil {
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

func (app *application) refreshFeed(w http.ResponseWriter, r *http.Request) {
	feedID, err := parseID(r)
	if err != nil {
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

	newEntries, err := syndication.GetNewEntries(feed.FeedUrl, feed.Type, feed.CheckedAt)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(newEntries) > 0 {
		now := time.Now().UTC().Format(time.RFC3339)
		entries := buildCreateEntryParams(feed.ID, now, newEntries)
		err = app.queries.CreateMultipleEntry(context.Background(), entries)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = app.queries.UpdateFeedCheckedAt(context.Background(), data.UpdateFeedCheckedAtParams{
			ID:        feed.ID,
			CheckedAt: now,
		})
		if err != nil {
			app.serverError(w, err)
			return
		}

	}
	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%d/", feedID))
	w.WriteHeader(http.StatusOK)
}

func (app *application) refreshAllFeeds(w http.ResponseWriter, r *http.Request) {
	err := app.refreshFeeds()
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Add("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (app *application) getEntry(w http.ResponseWriter, r *http.Request) {
	entryID, err := parseID(r)
	if err != nil {
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

	p := bluemonday.UGCPolicy()
	htmlContent := template.HTML(p.Sanitize(entry.Content))
	app.render(w, http.StatusOK, "entry.html", map[string]any{"entry": entry, "content": htmlContent})
}

func (app *application) deleteEntry(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.queries.DeleteEntry(context.Background(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	entryID, err := parseID(r)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.queries.MarkEntryRead(context.Background(), entryID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
