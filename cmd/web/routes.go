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
	// Legacy route redirects to canonical form.
	router.HandlerFunc(http.MethodGet, "/paste/view/:id", app.pasteLegacyRedirect)

	// Canonical paste view moved to /p/:slug to avoid conflict with /paste/create
	router.HandlerFunc(http.MethodGet, "/p/:slug", app.pasteView)
	router.HandlerFunc(http.MethodPost, "/p/:slug", app.pasteView)

	router.HandlerFunc(http.MethodGet, "/paste/create", app.pasteCreateForm)
	router.HandlerFunc(http.MethodPost, "/paste/create", app.pasteCreatePost)
	return router
}
