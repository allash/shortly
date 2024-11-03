package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
	
}