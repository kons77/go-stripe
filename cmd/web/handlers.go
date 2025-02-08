package main

import (
	"net/http"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["stripePubKey"] = app.config.stripe.key
	// app.infoLog.Println("public key", app.config.stripe.key)

	if err := app.renderTemplate(w, r, "terminal", &templateData{StringMap: stringMap}); err != nil {
		app.errorLog.Println(err)
	}
}
