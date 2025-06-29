package main

import (
	"net/http"
)

func (app *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/", app.health)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /feeds/", app.getFeeds)
	mux.HandleFunc("POST /feeds/", app.createFeed)
	mux.HandleFunc("GET /feeds/action/refresh-all/", app.refreshAllFeeds)
	mux.HandleFunc("GET /feeds/{id}/", app.getFeed)
	mux.HandleFunc("DELETE /feeds/{id}/", app.deleteFeed)
	mux.HandleFunc("POST /feeds/{id}/action/mark-read/", app.markFeedRead)
	mux.HandleFunc("GET /feeds/{id}/action/refresh/", app.refreshFeed)

	mux.HandleFunc("GET /entries/{id}/", app.getEntry)
	mux.HandleFunc("DELETE /entries/{id}/", app.deleteEntry)
	mux.HandleFunc("POST /entries/{id}/action/mark-read/", app.markEntryRead)
	mux.HandleFunc("POST /entries/action/mark-all-read/", app.markEntriesRead)

	return mux
}
