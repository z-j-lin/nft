package main

import (
	"log"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/monitor"
)

func main() {
	//connect to the redisDB
	rdb, err := redisDb.NewDBinstance()
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
	//start the que monitor
	qmon := monitor.NewQmon(rdb, eth)
	//start monitoring the MintQ list on redisDB
	qmon.StartTXQmon()
	//start the block monitor
	Bmon := monitor.NewBlockMon(eth)
	//start the monitor
	Bmon.Startmon()

}
