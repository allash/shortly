package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/:shortUrl", GetUrl)
	router.POST("/url/shorten", GenerateShortUrl)

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Access-Control-Request-Method") != "" {
			header := req.Header
			header.Set("Access-Control-Allow-Method", header.Get("Allow"))
			header.Set("Access-Controll-Allow-Origin", "*")
		}

		w.WriteHeader(http.StatusNoContent)
	})


	log.Fatal(http.ListenAndServe(":3000", router))
}

