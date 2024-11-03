package main

import (
	"net/http"
)

func (app *Application) logError(_ *http.Request, err error) {
	app.Logger.Println(err)
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	err := app.writeJSON(w, status, message, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "Internal server error"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "Resource not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}
