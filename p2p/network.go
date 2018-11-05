package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	b "golang-poc/blockchain"
	"io"
	"log"
	mrand "math/rand"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type Blockchain struct {
	Blocks []b.Block
}

func MakeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	r := createReaderWithRandomNumbers(randseed)

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	//if !secio {
	//	opts = append(opts, libp2p.NoEncryption())
	//}

	basicHost, err := libp2p.New(context.Background(), opts...)

	if err != nil {
		return nil, err
	}

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	return basicHost, nil
}

func CreateGenesisNode(s net.Stream, genesisBlock b.Block) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go handleStream(rw, genesisBlock)
}

func createReaderWithRandomNumbers(randseed int64) io.Reader {
	var r io.Reader
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	return r
}

func handleStream(rw *bufio.ReadWriter, genesisBlock b.Block) {
	blocks := append(make([]b.Block, 0), genesisBlock)
	blockchain := Blockchain{Blocks: blocks}
	blockchainChannel := make(chan Blockchain)
	blockchainChannel <- blockchain
	broadcastState(blockchainChannel, rw)
	for {
		str, err := rw.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}

		if str != "\n" {
			receivedChain := marshallReceivedChain(str)

			mutex.Lock()
			if len(receivedChain) > len(blockchain.Blocks) {
				blockchain.Blocks = receivedChain
				printBlockchain(blockchain.Blocks)
				blockchainChannel <- blockchain
			}
			mutex.Unlock()
		}
	}
}

func printBlockchain(Blockchain []b.Block) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
}

func marshallReceivedChain(str string) []b.Block {
	receivedChain := make([]b.Block, 0)
	if err := json.Unmarshal([]byte(str), &receivedChain); err != nil {
		log.Fatal(err)
	}
	return receivedChain
}

func broadcastState(blockchainChannel chan Blockchain, rw *bufio.ReadWriter) {
	currentBlockchain := Blockchain{Blocks: make([]b.Block, 0)}
	for {
		time.Sleep(5 * time.Second)
		select {
		case newBlockChain := <-blockchainChannel:
			currentBlockchain = newBlockChain
			broadcast(rw, marshallBlockchainToBytes(currentBlockchain))
		default:
			if len(currentBlockchain.Blocks) > 0 {
				broadcast(rw, marshallBlockchainToBytes(currentBlockchain))
			}
		}
	}
}

func marshallBlockchainToBytes(blockchain Blockchain) []byte {
	bytes, err := json.Marshal(blockchain.Blocks)
	if err != nil {
		log.Println(err)
	}
	return bytes
}

func broadcast(rw *bufio.ReadWriter, bytes []byte) {
	mutex.Lock()
	rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	rw.Flush()
	mutex.Unlock()
}
