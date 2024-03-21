package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	err := app.templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error executing template",
			"function": "home",
		})

		return
	}
}