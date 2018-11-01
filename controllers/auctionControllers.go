package controllers

import (
	"encoding/json"
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
		u.RespondWithError(w, "Error while decoding request body", http.StatusBadRequest)
		return
	}

	auction.AccountID = userId

	if resp, ok := auction.Validate(); !ok {
		u.RespondWithError(w, resp, http.StatusBadRequest)
		return
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
	auctionId, _ := strconv.ParseUint(u.GetPathVar("id", r), 10, 32)

	data := dao.GetAuction(uint(auctionId))

	u.Respond(w, data)
}
