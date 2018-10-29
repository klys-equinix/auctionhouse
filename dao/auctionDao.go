package dao

import (
	"github.com/jinzhu/gorm"
	u "golang-poc/utils"
	"log"
)

type Auction struct {
	CommonModelFields
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	AskingPrice  uint64        `json:"askingPrice"`
	AccountID    uint          `json:"accountId"`
	AuctionFiles []AuctionFile `json:"auctionFiles"`
}

func (auction *Auction) Validate() (map[string]interface{}, bool) {

	if auction.Name == "" {
		return u.Message(400, "Auction name should be on the payload"), false
	}

	if auction.Description == "" {
		return u.Message(400, "Description should be on the payload"), false
	}

	if auction.AccountID == 0 {
		return u.Message(400, "User is not recognized"), false
	}

	return u.Message(200, "success"), true
}

func (auction *Auction) Create() map[string]interface{} {

	if resp, ok := auction.Validate(); !ok {
		return resp
	}

	GetDB().Create(auction)

	resp := u.Message(201, "success")
	resp["auction"] = auction
	return resp
}

func GetAuction(id uint) *Auction {

	auction := &Auction{}
	err := GetDB().Table("auctions").Preload("AuctionFiles").Where("id = ?", id).First(auction).Error
	if err != nil {
		return nil
	}
	return auction
}

func GetAllAuctions(accountId uint64, name string) []*Auction {

	auctions := make([]*Auction, 0)
	db := GetDB().Table("auctions")

	db = buildAuctionQuery(accountId, name, db)

	err := db.Preload("AuctionFiles").Find(&auctions).Error

	if err != nil {
		log.Println(err)
		return nil
	}

	return auctions
}

func buildAuctionQuery(userId uint64, name string, db *gorm.DB) *gorm.DB {
	if userId != 0 {
		db = db.Where("account_id = ?", userId)
	}
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	return db
}
