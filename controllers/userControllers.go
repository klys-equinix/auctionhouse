package controllers

import (
	"golang-poc/dao"
	"golang-poc/dto"
	u "golang-poc/utils"
	"net/http"
)

var GetCurrentUser = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)

	data, err := dao.GetUser(uint(userId))

	if err != nil {
		u.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Respond(w, &dto.AccountDto{Email: data.Email, ID: data.ID})
}
