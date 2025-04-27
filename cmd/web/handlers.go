package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
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

	files := []string{
		"./ui/templates/base.html",
		"./ui/templates/partials/header.html",
		"./ui/templates/partials/entry_item.html",
		"./ui/templates/home.html",
	}

	ts, err := template.New("home").Funcs(template.FuncMap{
		"formatDate": formatDate,
	}).ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", unreadEntries)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.queries.GetFeeds(context.Background())
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/templates/base.html",
		"./ui/templates/partials/header.html",
		"./ui/templates/partials/feed_item.html",
		"./ui/templates/feeds.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// TODO: count unnecessary here
	err = ts.ExecuteTemplate(w, "base", map[string]any{
		"feeds": feeds,
		"count": len(feeds),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) getFeed(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedId < 1 {
		app.notFound(w)
		return
	}

	feed, err := app.queries.GetFeed(context.Background(), feedId)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	entries, err := app.queries.GetFeedEntries(context.Background(), feedId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/templates/base.html",
		"./ui/templates/partials/header.html",
		"./ui/templates/partials/entry_item.html",
		"./ui/templates/feed.html",
	}

	ts, err := template.New("feed").Funcs(template.FuncMap{
		"formatDate": formatDate},
	).ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	ts.ExecuteTemplate(w, "base", map[string]any{
		"feed":    feed,
		"entries": entries,
	})

}

func (app *application) createFeed(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	url := r.PostForm.Get("feedUrl")
	details, err := syndication.ExtractFeedDetails(url)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var feedURL, siteURL string
	for _, link := range details.Links {
		switch link.Rel {
		case "self":
			feedURL = link.Href
		default:
			siteURL = link.Href
		}
	}

	feed, err := app.queries.CreateFeed(context.Background(), data.CreateFeedParams{
		Title:     details.Title,
		Subtitle:  sql.NullString{String: details.Subtitle, Valid: details.Subtitle != ""},
		FeedUrl:   feedURL,
		SiteUrl:   siteURL,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			app.clientError(w, http.StatusConflict)
			return
		}
		app.serverError(w, err)
		return
	}

	for _, entry := range details.Entries {
		_, err := app.queries.CreateEntry(context.Background(), data.CreateEntryParams{
			FeedID:      feed.ID,
			Title:       entry.Title,
			Subtitle:    sql.NullString{String: entry.Subtitle, Valid: entry.Subtitle != ""},
			Author:      sql.NullString{String: entry.Author.Name, Valid: entry.Author.Name != ""},
			Content:     entry.Content,
			ExternalUrl: entry.Link.Href,
			PublishedAt: entry.Published,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		})

		if err != nil {
			app.serverError(w, err)
			return
		}

	}
	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%d/", feed.ID))
	w.WriteHeader(http.StatusCreated)
}

func (app *application) deleteFeed(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedId < 1 {
		app.notFound(w)
		return
	}

	err = app.queries.DeleteFeed(context.Background(), feedId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Add("HX-Redirect", "/feeds/")
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) getEntry(w http.ResponseWriter, r *http.Request) {
	entryId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || entryId < 1 {
		app.notFound(w)
		return
	}

	entry, err := app.queries.GetEntry(context.Background(), entryId)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	files := []string{
		"./ui/templates/base.html",
		"./ui/templates/partials/header.html",
		"./ui/templates/entry.html",
	}

	ts, err := template.New("entry").Funcs(template.FuncMap{
		"formatDate": formatDate,
		"ytNoCookie": ytNoCookie,
	}).ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	ts.ExecuteTemplate(w, "base", map[string]any{
		"entry":   entry,
		"content": template.HTML(entry.Content),
	})
}

func (app *application) markEntryRead(w http.ResponseWriter, r *http.Request) {
	entryId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || entryId < 1 {
		app.notFound(w)
		return
	}

	err = app.queries.UpdateEntryReadStatus(context.Background(), data.UpdateEntryReadStatusParams{
		Read: 1,
		ID:   entryId,
	})
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) markFeedRead(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedId < 1 {
		app.notFound(w)
		return
	}

	err = app.queries.MarkFeedRead(context.Background(), feedId)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) refreshFeed(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || feedId < 1 {
		app.notFound(w)
		return
	}

	feed, err := app.queries.GetFeed(context.Background(), feedId)
	if err != nil {
		app.notFound(w)
		return
	}

	newEntries, err := syndication.GetNewEntries(feed)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// TODO: Wrap this in a transaction
	now := time.Now().UTC().Format(time.RFC3339)
	for _, entry := range newEntries {
		_, err := app.queries.CreateEntry(context.Background(), data.CreateEntryParams{
			FeedID:      feed.ID,
			Title:       entry.Title,
			Subtitle:    sql.NullString{String: entry.Subtitle, Valid: entry.Subtitle != ""},
			Author:      sql.NullString{String: entry.Author.Name, Valid: entry.Author.Name != ""},
			Content:     entry.Content,
			ExternalUrl: entry.Link.Href,
			PublishedAt: entry.Published,
			CreatedAt:   now,
		})

		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	err = app.queries.UpdateFeedCheckedAt(context.Background(), data.UpdateFeedCheckedAtParams{
		ID:        feed.ID,
		CheckedAt: now,
	})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%d/", feed.ID))
	w.WriteHeader(http.StatusOK)
}
