package main

import (
	"log"

	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/monitor"
)

func main() {
	//connect to the redisDB
	rdb, err := monitor.NewDBinstance()
	if err != nil {
		log.Fatal(err)
	}
	//instantiate ether instance
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	contractAddr := "0x3DA85558aF6D0d0D03283fa23eD1edE90f7E3E03"
	eth, err := blockchain.NewEtherClient(rpcurl, contractAddr)
	if err != nil {
		log.Panic(err)
	}
	mintq := make(chan *blockchain.Transaction)
	//run loop for monitoring redisdb
	for {
		//check if there is something on the mintq
		account, resourceID, err := rdb.DQmint()
		// if there is a transaction
		if account != "" {
			//instantiate a new transaction
			tx := blockchain.NewTransaction(account)
			//put the transaction in the buffered channel
			//auto blocking when channel is full
			mintq <- tx
		}

		//run the transaction processor
		// if nothing is in the mintQ kill the processor
	}
}
