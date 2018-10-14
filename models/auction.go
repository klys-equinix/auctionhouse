package models

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
	UserId      uint   `json:"user_id"` //The user that this auction belongs to
}

func (auction *Auction) Validate() (map[string]interface{}, bool) {

	if auction.Name == "" {
		return u.Message(false, "Auction name should be on the payload"), false
	}

	if auction.Description == "" {
		return u.Message(false, "Description should be on the payload"), false
	}

	if auction.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	return u.Message(true, "success"), true
}

func (auction *Auction) Create() map[string]interface{} {

	if resp, ok := auction.Validate(); !ok {
		return resp
	}

	GetDB().Create(auction)

	resp := u.Message(true, "success")
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
	err := GetDB().Table("auctions").Where("user_id = ?", user).Find(&auctions).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return auctions
}
