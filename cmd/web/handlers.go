package main

import (
	"context"
	"net/http"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	version, err := app.queries.GetSQLiteVersion(context.Background())
	if err != nil {
		app.serverError(w, err)
	}
	w.Write([]byte(version))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html", nil)
}
