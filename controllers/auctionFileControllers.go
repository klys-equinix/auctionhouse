package controllers

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"golang-poc/dao"
	u "golang-poc/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var CreateAuctionFile = func(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	buf, name, extension := readFile(r, buf)

	vars := mux.Vars(r)
	auctionId, _ := strconv.ParseUint(vars["id"], 10, 32)
	auctionFile := &dao.AuctionFile{Name: name, Extension: extension, AuctionID: uint(auctionId)}

	if resp, ok := auctionFile.Validate(); !ok {
		u.RespondWithMessage(w, resp)
	}

	created, err := auctionFile.Create(buf)

	if err != nil {
		u.RespondWithError(w, u.Message(500, err.Error()), http.StatusInternalServerError)
		return
	}

	buf.Reset()

	u.Respond(w, created)
}

var GetAuctionFileById = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auctionFileId, _ := strconv.ParseUint(vars["fileId"], 10, 32)

	data, auctionFile := dao.GetAuctionFile(uint(auctionFileId))

	fileName := fmt.Sprintf("%s.%s", auctionFile.Name, auctionFile.Extension)

	u.RespondWithFile(w, data, fileName)
}

func readFile(r *http.Request, buf *bytes.Buffer) (*bytes.Buffer, string, string) {
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	io.Copy(buf, file)
	return buf, name[0], name[1]
}
