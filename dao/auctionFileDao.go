package dao

import (
	"bytes"
	"fmt"
	u "golang-poc/utils"
	"log"
	"os"
)

type AuctionFile struct {
	CommonModelFields
	Name      string `json:"name"`
	Extension string `json:"extension"`
	AuctionID uint   `json:"auctionId"`
}

func (auctionFile *AuctionFile) Validate() (map[string]interface{}, bool) {

	if auctionFile.Name == "" {
		return u.Message(400, "Auction name should be on the payload"), false
	}

	if auctionFile.Extension == "" {
		return u.Message(400, "Extension should be on the payload"), false
	}

	if auctionFile.AuctionID == 0 {
		return u.Message(400, "AuctionId should be on the payload"), false
	}

	return u.Message(200, "success"), true
}

func (auctionFile *AuctionFile) Create(buf *bytes.Buffer) (*AuctionFile, error) {
	tx := GetDB().Begin()

	if err := tx.Create(auctionFile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	success, err := auctionFile.SaveFile(buf)

	if !success {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return auctionFile, nil
}

func GetAuctionFile(id uint) (*os.File, *AuctionFile, error) {

	auctionFile := &AuctionFile{}
	err := GetDB().Table("auction_files").Where("id = ?", id).First(auctionFile).Error

	if err != nil {
		return nil, nil, err
	}

	filePath := auctionFile.buildFilePath()

	f, err := os.Open(filePath)

	if err != nil {
		return nil, nil, err
	}

	return f, auctionFile, nil
}

func (auctionFile *AuctionFile) SaveFile(buf *bytes.Buffer) (bool, error) {
	u.CreateStorageDirectoryIfNotExists()

	filePath := auctionFile.buildFilePath()

	if exists, err := u.FileExists(filePath); exists {
		return false, err
	}

	f, err := os.Create(filePath)

	if err != nil {
		log.Panicf("Cannot create file %s", err)
	}

	defer f.Close()
	_, err = f.Write(buf.Bytes())

	return err == nil, err
}

func (auctionFile *AuctionFile) buildFilePath() string {
	return fmt.Sprintf("%s/%d.%s", os.Getenv("storage_path"), int(auctionFile.ID), auctionFile.Extension)
}
