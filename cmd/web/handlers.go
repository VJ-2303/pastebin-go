package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vj-2303/pastebin-go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Pastebin")
}
func (app *application) pasteCreateForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "create.page.html", nil)
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
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	content := r.PostForm.Get("content")
	expiresStr := r.PostForm.Get("expires")
	password := r.PostForm.Get("password")

	if content == "" {
		fmt.Fprintln(w, "Content cannot be blank.")
		return
	}
	expires, err := strconv.Atoi(expiresStr)
	if err != nil || (expires != 4 && expires != 24) {
		fmt.Fprintln(w, "Invalid expiry value.")
		return
	}
	uniqueString, err := app.generateRandomString(8)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var passwordHash []byte
	if password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(passwordHash), 12)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	expiryTime := time.Now().Add(time.Duration(expires) * time.Hour)

	id, err := app.pastes.Insert(uniqueString, content, passwordHash, expiryTime)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/paste/view/%d", id), http.StatusSeeOther)
}
