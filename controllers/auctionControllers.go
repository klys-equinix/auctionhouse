package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang-poc/dao"
	u "golang-poc/utils"
	"net/http"
	"strconv"
)

var CreateAuction = func(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("user").(uint)
	auction := &dao.Auction{}

	err := json.NewDecoder(r.Body).Decode(auction)

	if err != nil {
		u.RespondWithError(w, u.Message(400, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	auction.AccountID = userId

	if resp, ok := auction.Validate(); !ok {
		u.RespondWithMessage(w, resp)
	}

	u.Respond(w, auction.Create())
}

var GetAllAuctions = func(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	accountId, _ := strconv.ParseUint(v.Get("accountId"), 10, 32)

	data := dao.GetAllAuctions(accountId, v.Get("name"))

	u.Respond(w, data)
}

var GetAuctionById = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auctionId, _ := strconv.ParseUint(vars["id"], 10, 32)

	data := dao.GetAuction(uint(auctionId))

	u.Respond(w, data)
}
