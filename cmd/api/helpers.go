package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (app *Application) readJSON(_ http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("malformed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}