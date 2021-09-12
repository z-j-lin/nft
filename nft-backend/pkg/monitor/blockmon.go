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
	state      objects.State
	taskClient *tasks.TaskClient
}

//makes an object to start a block monitor from
//I might make this a go routine that also call Startmon
func NewBlockMon(ether *blockchain.Ethereum) *monitor {
	redisAddr := "127.0.0.1:6379"
	TC := tasks.NewTaskClient(redisAddr)
	rdb, err := redisDb.NewDBinstance()
	if err != nil {
		log.Fatalf("failed to connect with redisdb: %v", err)
	}
	return &monitor{
		eth:        ether,
		db:         rdb,
		taskClient: TC,
	}
}

func (mon *monitor) Startmon() <-chan bool {
	//find the initial state in db

	initState, err := mon.db.GetState()
	mon.state = initState
	//if initial state is not in db, initiate a state starting from
	//initial block of contract
	if err != nil {
		RootBlock := uint64(10910043)
		initState = objects.State{
			HighestFinalizedBlock: RootBlock,
			HighestProcessedBlock: RootBlock,
		}
		//record the state
		mon.state = initState
	}
	//might not need this
	killChan := make(chan bool)
	//starts the monitoring loop go routine
	go mon.monitorloop(mon.state, killChan)
	//DO i need this?
	return killChan
}

//this function querys for a block every 5 seconds
func (mon *monitor) monitorloop(state objects.State, exit chan bool) error {
	run := make(chan bool)
	for {
		if !state.InSync {
			run <- true
		}
		select {
		case <-exit:
			log.Println("recieved kill signal")
			return nil
		case <-time.After(5 * time.Second):
			mon.getBlock(&state)
		case <-run:
			mon.getBlock(&state)
		}
	}
}
func (mon *monitor) getBlock(state *objects.State) {
	//get latest block
	header, err := mon.eth.Client.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		log.Fatalf("at monitorloop %v", err)
	}
	//Most recent block added to the block chain
	latestBlock := header.Number.Uint64()
	//40 blocks below the most recent block
	delayedLatestBlock := latestBlock - uint64(40)

	//next block to process
	currentBlock := state.HighestProcessedBlock + 1
	if delayedLatestBlock < currentBlock+20 {
		state.InSync = true
	}
	if currentBlock < delayedLatestBlock {
		mon.taskClient.QVerificationTask(int64(currentBlock))
		state.HighestProcessedBlock = currentBlock
		//update state on redis
		err = mon.db.UpdateState(state)
		if err != nil {
			log.Fatal("unable to update state", err)
		}
	}
}
