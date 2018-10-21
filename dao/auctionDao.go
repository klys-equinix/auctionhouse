package dao

import (
	u "../utils"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Auction struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	AskingPrice uint64 `json:"askingPrice"`
	AccountId   uint   `gorm:"type:bigint REFERENCES accounts(id)"`
}

func (auction *Auction) Validate() (map[string]interface{}, bool) {

	if auction.Name == "" {
		return u.Message(400, "Auction name should be on the payload"), false
	}

	if auction.Description == "" {
		return u.Message(400, "Description should be on the payload"), false
	}

	if auction.AccountId == 0 {
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
	err := GetDB().Table("auctions").Where("id = ?", id).First(auction).Error
	if err != nil {
		return nil
	}
	return auction
}

func GetAuctions(user uint) []*Auction {

	auctions := make([]*Auction, 0)
	err := GetDB().Table("auctions").Where("account_id = ?", user).Find(&auctions).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return auctions
}
