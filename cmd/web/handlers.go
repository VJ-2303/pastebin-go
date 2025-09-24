package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vj-2303/pastebin-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Pastebin")
}
func (app *application) pasteCreateForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Displaying the form for creating a new paste...")
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	paste, err := app.pastes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := &templateData{
		Paste: paste,
	}
	app.render(w, http.StatusOK, "view.page.html", data)
}

func (app *application) pasteCreatePost(w http.ResponseWriter, r *http.Request) {

	uniqueString := "adcdef"
	content := "This is the content of the new paste"
	expires := time.Now().Add(24 * time.Hour)

	id, err := app.pastes.Insert(uniqueString, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/paste/view/%d", id), http.StatusSeeOther)
}
