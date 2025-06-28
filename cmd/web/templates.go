package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var functions = template.FuncMap{
	"formatDate": formatDate,
	"ytNoCookie": ytNoCookie,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/templates/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/templates/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func (app *application) render(w http.ResponseWriter, status int, page string, data any) {
	ts, ok := app.templates[page]
	if !ok {
		err := fmt.Errorf("the template '%s' does not exist", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func formatDate(dt string) string {
	val, _ := time.Parse(time.RFC3339, dt)
	return val.Format("Jan 02, 2006")
}

func ytNoCookie(ytURL string) string {
	_, videoID, _ := strings.Cut(ytURL, "watch?v=")
	return fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s", videoID)
}
