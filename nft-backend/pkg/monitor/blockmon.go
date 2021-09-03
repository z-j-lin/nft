package monitor

import (
	"context"
	"log"
	"time"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	objects "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Objects"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type monitor struct {
	eth   *blockchain.Ethereum
	db    redisDb.Database
	state *objects.State
}

func (mon *monitor) startmon() <-chan bool {
	//find the initial state in db
	//if initial state is not in db
	initState, err := mon.db.GetState()
	if err != nil {
		RootBlock := uint64(10910043)
		initState = &objects.State{
			HighestFinalizedBlock: RootBlock,
			HighestProcessedBlock: RootBlock,
		}
	}
	//might not need this
	killChan := make(chan bool)
	go mon.monitorloop(initState, killChan)
	//DO i need this?
	return killChan
}

func (mon *monitor) monitorloop(state *objects.State, exit <-chan bool) error {
	for {
		select {
		case <-exit:
			log.Println("recieved kill signal")
			return nil
		case <-time.After(5 * time.Second):
			//get latest block
			header, err := mon.eth.Client.HeaderByNumber(context.TODO(), nil)
			if err != nil {
				log.Fatalf("at monitorloop %v", err)
			}
			//highest block added to verification queue
			currentBlock := state.HighestProcessedBlock + 1
			//Most recent block added to the block chain
			latestBlock := header.Number.Uint64()
			//40 blocks below the most recent block
			delayedLatestBlock := latestBlock - uint64(40)
			/*add a block to the verfication queue if
			the currentblock is less than or eq to delayedLatestBlock */
			if currentBlock <= delayedLatestBlock {
				//queue the pendingblock
				mon.db.QPendingBlock(currentBlock)
				state.HighestProcessedBlock = currentBlock
				//update state on "disk"
				mon.db.updateState(state * State)
			}
		}
	}
}
