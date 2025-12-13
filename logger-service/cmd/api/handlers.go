package main

import (
	"logger-service/data"
	"net/http"
)

type requestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.Logs.Insert(event)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged successfully",
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}
