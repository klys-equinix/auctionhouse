package dao

import (
	u "../utils"
	"bytes"
	"fmt"
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

func (auctionFile *AuctionFile) Create(buf *bytes.Buffer) map[string]interface{} {
	tx := GetDB().Begin()

	if resp, ok := auctionFile.Validate(); !ok {
		return resp
	}

	if err := tx.Create(auctionFile).Error; err != nil {
		tx.Rollback()
		return u.Message(500, err.Error())
	}

	success, err := auctionFile.SaveFile(buf)

	if !success {
		tx.Rollback()
		return u.Message(500, err.Error())
	}

	tx.Commit()

	resp := u.Message(201, "success")
	resp["auctionFile"] = auctionFile
	return resp
}

func (auctionFile *AuctionFile) SaveFile(buf *bytes.Buffer) (bool, error) {
	createStorageDirectoryIfNotExists()

	filePath := fmt.Sprintf("%s/%d.%s", os.Getenv("storage_path"), int(auctionFile.ID), auctionFile.Extension)

	if exists, err := exists(filePath); exists {
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

func GetAuctionFile(id uint) *AuctionFile {

	auctionFile := &AuctionFile{}
	err := GetDB().Table("auction_files").Where("id = ?", id).First(auctionFile).Error
	if err != nil {
		return nil
	}
	return auctionFile
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

func createStorageDirectoryIfNotExists() {
	storagePath := os.Getenv("storage_path")
	if exists, err := exists(storagePath); !exists {
		err = os.MkdirAll(storagePath, 0755)
		if err != nil {
			panic(err)
		}
	}
}
