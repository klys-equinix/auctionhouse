package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Offer struct {
	BuyerId       uint
	OfferedAmount int
}

type Block struct {
	Index           int
	MinimalOffer    uint64
	Timestamp       string
	TerminationTime string
	Hash            string
	PrevHash        string
	Offer           Offer
	SellerId        uint
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	if oldBlock.Offer.OfferedAmount >= newBlock.Offer.OfferedAmount {
		return false
	}

	if newBlock.TerminationTime == "" {
		return false
	}

	return true
}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Offer.OfferedAmount) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateBlock(oldBlock Block, offeredAmount int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Offer.OfferedAmount = offeredAmount
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func GenerateGenesisBlock(initialOffer int, sellerId uint, auctionEnd string, askingPrice uint64) Block {

	newBlock := Block{}
	genesisOffer := Offer{OfferedAmount: initialOffer}
	t := time.Now()

	newBlock.Index = 0
	newBlock.Timestamp = t.String()
	newBlock.Offer = genesisOffer
	newBlock.TerminationTime = auctionEnd
	newBlock.SellerId = sellerId
	newBlock.Hash = calculateHash(newBlock)
	newBlock.MinimalOffer = askingPrice

	return newBlock
}
