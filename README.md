# golang-poc
Proof of concept -> creation of REST API with GoLang

run -> go build main.go

database credentials, storage path and configuration -> .env

example database -> docker-compose.yml

basic requestst -> auction_house.postman_collection.json

THE PAIN AND SUFFERING - installing libp2p, because those fuckers have to use a decentralised package manager (gx based on ipfs)

1. go get -u -d github.com/libp2p/go-libp2p/... (MAKE SURE THAT GOPATH is set to global GOPATH, AND DO THAT IN $GOPATH)
2. cd $GOPATH/src/github.com/libp2p/go-libp2p
3. make
4. make deps
