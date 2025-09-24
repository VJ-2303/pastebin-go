package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/paste/view/:id", app.pasteView)

	router.HandlerFunc(http.MethodGet, "/paste/create", app.pasteCreateForm)
	router.HandlerFunc(http.MethodPost, "/paste/create", app.pasteCreatePost)
	return router
}
