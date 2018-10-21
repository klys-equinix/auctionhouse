package dao

import (
	u "../utils"
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
)

type AuctionFile struct {
	gorm.Model
	Name      string `json:"name"`
	Extension string `json:"extension"`
	AuctionId uint   `gorm:"type:bigint REFERENCES auctions(id)"`
}

func (auctionFile *AuctionFile) Validate() (map[string]interface{}, bool) {

	if auctionFile.Name == "" {
		return u.Message(400, "Auction name should be on the payload"), false
	}

	if auctionFile.Extension == "" {
		return u.Message(400, "Extension should be on the payload"), false
	}

	if auctionFile.AuctionId == 0 {
		return u.Message(400, "AuctionId should be on the payload"), false
	}

	return u.Message(200, "success"), true
}

func (auctionFile *AuctionFile) Create(buf *bytes.Buffer) map[string]interface{} {
	tx := GetDB().Begin()

	if resp, ok := auctionFile.Validate(); !ok {
		return resp
	}

	tx.Create(auctionFile)
	success, err := auctionFile.SaveFile(buf)
	if !success {
		tx.Rollback()
		resp := u.Message(500, err.Error())
		return resp
	}

	tx.Commit()

	resp := u.Message(201, "success")
	resp["auctionFile"] = auctionFile
	return resp
}

func (auctionFile *AuctionFile) SaveFile(buf *bytes.Buffer) (bool, error) {
	filePath := strconv.Itoa(int(auctionFile.ID)) + "." + auctionFile.Extension
	if exists, err := exists(filePath); exists {
		return false, err
	}
	f, err := os.Create(filePath)
	defer f.Close()
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return false, err
	}
	return true, err
}

func GetAuctionFile(id uint) *AuctionFile {

	auctionFile := &AuctionFile{}
	err := GetDB().Table("auction_files").Where("id = ?", id).First(auctionFile).Error
	if err != nil {
		return nil
	}
	return auctionFile
}

func GetAuctionFilesForAuction(user uint) []*AuctionFile {

	auctionFiles := make([]*AuctionFile, 0)
	err := GetDB().Table("auction_files").Where("auction_id = ?", user).Find(&auctionFiles).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return auctionFiles
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
