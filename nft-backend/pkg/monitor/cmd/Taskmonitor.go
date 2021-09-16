package main

import (
	"log"
	"math/big"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/monitor"
)

func kill(killchan chan bool) { killchan <- true }
func main() {
	//connect to the redisDB
	rdb, err := redisDb.NewDBinstance()
	if err != nil {
		log.Fatal(err)
	}
	//instantiate ether instance
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	contractAddr := "0xb410756d52b1250aB9bE358437Ab41a4D7636Af8"
	eth, err := blockchain.NewEtherClient(rpcurl, contractAddr, big.NewInt(int64(3)))
	if err != nil {
		log.Panic(err)
	}
	//start the block monitor
	Bmon := monitor.NewBlockMon(eth)
	log.Println("Block Monitor Started")
	//enqueues verification tasks onto the task queue
	go Bmon.Startmon()
	defer kill(Bmon.Kill)
	//start the task que
	// handles validation and tx tasks
	log.Println("Task server Started")
	monitor.NewQmon(rdb.Client.Options().Addr, 0, eth, rdb)
}
