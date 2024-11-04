package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/health", app.health)
	router.GET("/v1/:shortUrl", app.getLongUrl)
	router.POST("/v1/url/shorten", app.createShortUrl)

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Access-Control-Request-Method") != "" {
			header := req.Header
			header.Set("Access-Control-Allow-Method", header.Get("Allow"))
			header.Set("Access-Controll-Allow-Origin", "*")
		}

		w.WriteHeader(http.StatusNoContent)
	})

	return router
}