package controllers

import (
	"../dao"
	u "../utils"
	"encoding/json"
	"net/http"
)

var CreateAuction = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	auction := &dao.Auction{}

	err := json.NewDecoder(r.Body).Decode(auction)
	if err != nil {
		u.Respond(w, u.Message(400, "Error while decoding request body"))
		return
	}

	auction.UserId = user
	resp := auction.Create()
	u.Respond(w, resp)
}

var GetAuctionsFor = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user").(uint)
	data := dao.GetAuctions(id)
	resp := u.Message(200, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
