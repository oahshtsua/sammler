package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

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
