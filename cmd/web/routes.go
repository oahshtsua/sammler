package main

import (
	"net/http"
)

func (app *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/", app.health)

	mux.HandleFunc("GET /", app.home)

	return mux
}
