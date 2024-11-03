package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var urlMappings = make(map[string]string)

func GenerateShortUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("Received request \n")

	var payload LongUrl

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	urlMappings["123"] = payload.Value

	response := &ShortUrlResponse{Value: "123"}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func GetUrl(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	shortUrl := params.ByName("shortUrl")
	fmt.Printf("Params: %s \n", params.ByName("shortUrl"))

	longUrl := urlMappings[shortUrl]

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(longUrl)
}