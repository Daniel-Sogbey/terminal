package main

import (
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)

	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "Invalid request body"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	//TODO: authenticate
	app.infoLog.Printf("%+v", creds)

	//TODO: send back response
	payload.Error = false
	payload.Message = "Signed in"

	err = app.writeJSON(w, http.StatusOK, payload)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

}
