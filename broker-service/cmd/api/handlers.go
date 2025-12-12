package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   authPayload `json:"auth,omitempty"`
}

type authPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w, requestPayload.Auth)
	default:
		app.ErrorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
		return
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, authPayload authPayload) {
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		app.ErrorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}
