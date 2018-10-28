package controllers

import (
	"../dao"
	u "../utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var CreateAuction = func(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("user").(uint)
	auction := &dao.Auction{}

	err := json.NewDecoder(r.Body).Decode(auction)
	if err != nil {
		u.Respond(w, u.Message(400, "Error while decoding request body"))
		return
	}

	auction.AccountID = userId
	resp := auction.Create()
	u.Respond(w, resp)
}

var GetAuctionsForUser = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user").(uint)
	data := dao.GetAuctionsForUser(id)
	resp := u.Message(200, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetAllAuctions = func(w http.ResponseWriter, r *http.Request) {
	data := dao.GetAllAuctions()
	resp := u.Message(200, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetAuctionById = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auctionId, _ := strconv.ParseUint(vars["id"], 10, 32)

	data := dao.GetAuction(uint(auctionId))
	resp := u.Message(200, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
