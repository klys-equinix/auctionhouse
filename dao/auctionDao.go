package dao

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Auction struct {
	CommonModelFields
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	AskingPrice     uint64        `json:"askingPrice"`
	AccountID       uint          `json:"accountId"`
	AuctionFiles    []AuctionFile `json:"auctionFiles"`
	TerminationTime time.Time     `json:"terminationTime"`
	AuctionHost     string        `json:"auctionHost"`
}

func (auction *Auction) Validate() (string, bool) {

	if auction.Name == "" {
		return "Auction name should be on the payload", false
	}

	if auction.Description == "" {
		return "Description should be on the payload", false
	}

	if auction.AccountID == 0 {
		return "User is not recognized", false
	}

	return "", true
}

func (auction *Auction) Create() *Auction {

	GetDB().Create(auction)

	return auction
}

func (auction *Auction) Start() *Auction {

	return auction
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
