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
	q := r.URL.Query().Get("q")
	if q != "" {
		paste, err := app.pastes.GetByUnique(q)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				pastes, lerr := app.pastes.Latest(20)
				if lerr != nil {
					app.serverError(w, lerr)
					return
				}
				data := &templateData{Pastes: pastes, SearchQuery: q, SearchError: "No paste found for that identifier"}
				app.render(w, http.StatusOK, "home.page.html", data)
				return
			}
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/p/%s", paste.UniqueString), http.StatusSeeOther)
		return
	}
	pastes, err := app.pastes.Latest(20)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Pastes: pastes}
	app.render(w, http.StatusOK, "home.page.html", data)
}
func (app *application) pasteCreateForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "create.page.html", nil)
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idParam := params.ByName("slug")

	var paste *models.Paste
	var err error
	if id, convErr := strconv.Atoi(idParam); convErr == nil {
		paste, err = app.pastes.Get(id)
	} else {
		paste, err = app.pastes.GetByUnique(idParam)
	}
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	if len(paste.PasswordHash) > 0 {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			provided := r.PostForm.Get("password")
			if bcrypt.CompareHashAndPassword(paste.PasswordHash, []byte(provided)) != nil {
				data := &templateData{Paste: nil, ViewError: "Incorrect password", PasswordForm: true, UniqueString: paste.UniqueString}
				app.render(w, http.StatusUnauthorized, "view.page.html", data)
				return
			}
		} else {
			data := &templateData{PasswordForm: true, UniqueString: paste.UniqueString}
			app.render(w, http.StatusOK, "view.page.html", data)
			return
		}
	}

	data := &templateData{Paste: paste, UniqueString: paste.UniqueString}
	app.render(w, http.StatusOK, "view.page.html", data)
}

func (app *application) pasteLegacyRedirect(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
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
	http.Redirect(w, r, fmt.Sprintf("/p/%s", paste.UniqueString), http.StatusMovedPermanently)
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
	uniqueString, err := app.generateRandomString(4)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var passwordHash []byte
	if password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	expiryTime := time.Now().Add(time.Duration(expires) * time.Hour)

	_, err = app.pastes.Insert(uniqueString, content, passwordHash, expiryTime)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/p/%s", uniqueString), http.StatusSeeOther)
}
