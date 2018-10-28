package controllers

import (
	"../dao"
	"../dto"
	u "../utils"
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	createAccountDto := &dto.CreateAccountDto{}
	err := json.NewDecoder(r.Body).Decode(createAccountDto)
	if err != nil {
		u.Respond(w, u.Message(400, "Invalid request"))
		return
	}

	if resp, ok := createAccountDto.Validate(); !ok {
		u.Respond(w, resp)
		return
	}

	resp := dao.Create(createAccountDto) //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &dao.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(400, "Invalid request"))
		return
	}

	resp := dao.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
