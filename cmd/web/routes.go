package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/paste/view/:id", app.pasteView)

	router.HandlerFunc(http.MethodGet, "/paste/create", app.pasteCreateForm)
	router.HandlerFunc(http.MethodPost, "/paste/create", app.pasteCreatePost)
	return router
}
