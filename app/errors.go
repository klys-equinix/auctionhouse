package app

import (
	u "golang-poc/utils"
	"net/http"
)

var NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	u.RespondWithMessage(w, u.Message(404, "This resources was not found on our server"))
})
