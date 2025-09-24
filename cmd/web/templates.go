package main

import (
	"html/template"
	"path/filepath"

	"github.com/vj-2303/pastebin-go/internal/models"
)

type templateData struct {
	Paste        *models.Paste
	PasswordForm bool
	ViewError    string
	UniqueString string
	Pastes       []*models.Paste
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.page.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles("./ui/html/layouts/base.layout.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.partial.html")
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
