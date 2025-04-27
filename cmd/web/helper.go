package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func formatDate(dt string) string {
	val, _ := time.Parse(time.RFC3339, dt)
	return val.Format("Jan 02, 2006")
}

func ytNoCookie(ytURL string) string {
	_, videoID, _ := strings.Cut(ytURL, "watch?v=")
	return fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s", videoID)
}

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

func (app *application) renderPartial(w http.ResponseWriter, name string, data any) {
	ts, err := template.ParseGlob("./ui/templates/partials/*.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, name, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

}
