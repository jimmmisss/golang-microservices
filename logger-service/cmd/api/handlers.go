package main

import (
	"log"
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)
	log.Println("Name:", requestPayload.Name)
	log.Println("Data:", requestPayload.Data)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	log.Println("Data 02:", event.Name)

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
