package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"shorturl.allash.com/internal/data"
)

var urlMappings = make(map[string]string)

func (app *Application) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := &data.HealthStatus {
		Status: "available",
		Environment: app.Config.Environment,
		Version: "1.0",
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) generateShortUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("Received request \n")

	var payload data.LongUrl

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	urlMappings["123"] = payload.Value

	response := &data.ShortUrlResponse{Value: "123"}
	err := app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getUrl(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	shortUrl := params.ByName("shortUrl")
	fmt.Printf("Params: %s \n", params.ByName("shortUrl"))

	longUrl := urlMappings[shortUrl]

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(longUrl)
}