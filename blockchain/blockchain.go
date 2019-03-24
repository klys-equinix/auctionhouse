package blockchain

type Blockchain struct {
	Blocks []Block
}

func (blockchain *Blockchain) IsBlockchainValid() bool {
	for i, block := range blockchain.Blocks {
		if i > 0 && !IsBlockValid(block, blockchain.Blocks[i-1]) {
			return false
		}
	}
	return true
}
