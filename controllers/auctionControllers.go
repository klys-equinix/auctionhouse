package controllers

import (
	"encoding/json"
	"fmt"
	b "golang-poc/blockchain"
	"golang-poc/dao"
	"golang-poc/p2p"
	u "golang-poc/utils"
	"gx/ipfs/QmYrWiWM4qtrnCeT3R14jY3ZZyirDNJgwK57q4qFYePgbd/go-libp2p-host"
	"net/http"
	"os"
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

	auction.Account.ID = userId

	if resp, ok := auction.Validate(); !ok {
		u.RespondWithError(w, resp, http.StatusBadRequest)
		return
	}

	p2pPort, _ := strconv.Atoi(os.Getenv("p2p_port"))

	auctionHost, addr, p2pErr := p2p.MakeBasicHost(p2pPort, false, 0)
	if p2pErr != nil {
		fmt.Print(p2pErr)
	}

	auction.AuctionHost = addr

	created := auction.Create()

	go CreateGenesisNode(auctionHost, created)

	u.Respond(w, created)
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

func CreateGenesisNode(h host.Host, auction *dao.Auction) {
	genesisBlock := b.GenerateGenesisBlock(0, auction.Account.ID, auction.TerminationTime.String(), auction.AskingPrice)
	streamHandler := p2p.GetStreamHandler(genesisBlock)
	h.SetStreamHandler("/p2p/1.0.0", streamHandler)
	select {}
}
