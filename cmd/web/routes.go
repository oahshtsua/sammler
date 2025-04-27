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

	return mux
}
