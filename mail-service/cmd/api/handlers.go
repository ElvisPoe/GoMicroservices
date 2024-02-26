package main

import (
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	type mailMessage struct {
		From     string `json:"from"`
		FromName string `json:"fromName"`
		To       string `json:"to"`
		Subject  string `json:"subject"`
		Message  string `json:"message"`
	}

	var requestPayload mailMessage

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	if err := app.Mailer.SendSMTPMessage(msg); err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Message sent to + " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
