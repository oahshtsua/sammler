package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

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

func parseID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("Invalid ID")
	}
	return id, nil
}

func buildCreateEntryParams(feedID int64, now string, entries []syndication.FeedEntry) []data.CreateEntryParams {
	params := make([]data.CreateEntryParams, 0, len(entries))
	for _, entry := range entries {
		params = append(params, data.CreateEntryParams{
			FeedID: feedID,
			Title:  entry.Title,
			Author: sql.NullString{
				String: entry.Author,
				Valid:  entry.Author != "",
			},
			Content:     entry.Content,
			ExternalUrl: entry.Link,
			PublishedAt: entry.Published,
			CreatedAt:   now,
		})

	}
	return params
}
