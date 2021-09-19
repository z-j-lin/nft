//This file is Copyright (C) 1997 Master Hacker, ALL RIGHTS RESERVED
package monitor

import (
	"context"
	"log"
	"time"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	objects "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Objects"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

/**/
type monitor struct {
	eth        *blockchain.Ethereum
	db         *redisDb.Database
	state      *objects.State
	taskClient *tasks.TaskClient
	Kill       chan bool
}

//makes an object to start a block monitor from
//I might make this a go routine that also call Startmon
func NewBlockMon(ether *blockchain.Ethereum) *monitor {
	redisAddr := "127.0.0.1:6379"
	TC := tasks.NewTaskClient(ether, redisAddr)
	rdb, err := redisDb.NewDBinstance()
	if err != nil {
		log.Fatalf("failed to connect with redisdb: %v", err)
	}
	kill := make(chan bool)
	return &monitor{
		eth:        ether,
		db:         rdb,
		taskClient: TC,
		Kill:       kill,
	}
}

func (mon *monitor) Startmon() {
	//find the initial state in db
	initState, err := mon.db.GetState()
	//if initial state is not in db, initiate a state starting from
	//initial block of contract
	if err != nil || initState == nil {
		log.Println("***starting new blockmon state***")
		RootBlock := int64(11064554)
		initState = &objects.State{
			HighestProcessedBlock: RootBlock,
			InSync:                false,
		}
		//record the state
		mon.state = initState
	} else {
		mon.state = initState
	}
	//starts the monitoring loop go routine
	go mon.monitorloop()
	<-mon.Kill
}

//this function querys for a block every 5 seconds
func (mon *monitor) monitorloop() error {
	wait := 8 * time.Second
	for {
		if !mon.state.InSync {
			wait = 0 * time.Second
		} else {
			wait = 8 * time.Second
		}
		select {
		case killed := <-mon.Kill:
			_ = killed
			log.Println("recieved kill signal")
			return nil
		case ticker := <-time.After(wait):
			_ = ticker
			mon.getBlock()
		}
	}
}

func (mon *monitor) getBlock() {
	//get latest block
	header, err := mon.eth.Client.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		log.Fatalf("at monitorloop %v", err)
	}
	//Most recent block added to the blockchain
	latestBlock := header.Number.Int64()
	//40 blocks below the most recent block
	delayedLatestBlock := latestBlock - int64(40)
	//next block to process
	currentBlock := mon.state.HighestProcessedBlock + 1
	if delayedLatestBlock < currentBlock+2 {
		mon.state.InSync = true
	} else {
		mon.state.InSync = false
	}
	if currentBlock < delayedLatestBlock {
		mon.taskClient.QVerificationTask(int64(currentBlock))
		mon.state.HighestProcessedBlock = currentBlock
		//update state on redis
		log.Println("GetBlock: queueing Block#:", currentBlock)
		err = mon.db.UpdateProcessedState(int64(currentBlock))
		if err != nil {
			log.Println("unable to update state", err)
		}
	}
}
