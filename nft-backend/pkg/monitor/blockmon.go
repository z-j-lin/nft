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
	Kill       chan bool
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
	mon.state = initState
	//if initial state is not in db, initiate a state starting from
	//initial block of contract
	if err != nil {
		log.Println("***starting new blockmon state***")
		RootBlock := uint64(11048465)
		initState = objects.State{
			HighestFinalizedBlock: RootBlock,
			HighestProcessedBlock: RootBlock,
			InSync:                false,
		}
		//record the state
		mon.state = initState
	}
	//starts the monitoring loop go routine
	log.Println("startomg monitor loop")
	go mon.monitorloop(mon.state, mon.Kill)
	<-mon.Kill
}

//this function querys for a block every 5 seconds
func (mon *monitor) monitorloop(state objects.State, exit <-chan bool) error {

	for {
		wait := 5 * time.Second
		if !state.InSync {
			wait = 0 * time.Second
		} else {
			wait = 5 * time.Second
		}
		log.Println("monitorloop running")
		select {
		case killed := <-exit:
			_ = killed
			log.Println("recieved kill signal")
			return nil
		case ticker := <-time.After(wait * time.Second):
			_ = ticker
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
	log.Println("currentBlock:", currentBlock, "delayedLatestBlock:", delayedLatestBlock)
	if delayedLatestBlock < currentBlock+20 {
		state.InSync = true
	} else {
		state.InSync = false
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
