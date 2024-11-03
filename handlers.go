package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetHello(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Printf("Received request \n")

	response := &ShortUrlResponse{}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}