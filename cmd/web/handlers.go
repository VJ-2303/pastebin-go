package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Pastebin")
}
func (app *application) pasteCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Displaying the form for creating a new paste...")
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Displaying a specific paste...")
}
