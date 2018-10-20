package app

import (
	u "../utils"
	"net/http"
)

var NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	u.Respond(w, u.Message(404, "This resources was not found on our server"))
})
