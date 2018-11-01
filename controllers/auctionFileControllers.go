package controllers

import (
	"bytes"
	"fmt"
	"golang-poc/dao"
	u "golang-poc/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var CreateAuctionFile = func(w http.ResponseWriter, r *http.Request) {
	buf, name, extension := readFile(r)

	auctionId, _ := strconv.ParseUint(u.GetPathVar("id", r), 10, 32)
	auctionFile := &dao.AuctionFile{Name: name, Extension: extension, AuctionID: uint(auctionId)}

	if resp, ok := auctionFile.Validate(); !ok {
		u.RespondWithMessage(w, resp)
		return
	}

	created, err := auctionFile.Create(buf)

	if err != nil {
		u.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.Reset()

	u.Respond(w, created)
}

var GetAuctionFileById = func(w http.ResponseWriter, r *http.Request) {
	auctionFileId, _ := strconv.ParseUint(u.GetPathVar("fileId", r), 10, 32)

	data, auctionFile := dao.GetAuctionFile(uint(auctionFileId))

	fileName := fmt.Sprintf("%s.%s", auctionFile.Name, auctionFile.Extension)

	u.RespondWithFile(w, data, fileName)
}

func readFile(r *http.Request) (*bytes.Buffer, string, string) {
	buf := &bytes.Buffer{}
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	io.Copy(buf, file)
	return buf, name[0], name[1]
}
