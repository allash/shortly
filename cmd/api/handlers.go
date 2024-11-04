package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"shortly.allash.com/internal/data"
	"shortly.allash.com/internal/generator"
)

func (app *Application) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := &data.HealthStatus {
		Status: "available",
		Environment: app.config.environment,
		Version: "1.0",
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) createShortUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var payload data.LongUrl

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	sf, err := generator.NewSnowflake(1, 1)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newId, err := sf.NextId()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	shortUrl := generator.Encode(newId)

	urlMapping := &data.UrlMapping{
		ID: newId,
		ShortUrl: shortUrl,
		LongUrl: payload.Value,
	}
	
	err = app.models.UrlMappings.Insert(urlMapping)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := &data.ShortUrlResponse{Value: shortUrl}
	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getLongUrl(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	shortUrl := params.ByName("shortUrl")

	longUrl, err := app.models.UrlMappings.Get(shortUrl)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	response := &data.LongUrl{Value: *longUrl}
	err = app.writeJSON(w, http.StatusOK, response, nil)
	if (err != nil) {
		app.serverErrorResponse(w, r, err)
	}
}