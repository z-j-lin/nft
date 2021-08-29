package monitor

import (
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type txqmon struct {
	db    Database
	MintQ chan *blockchain.transaction
}

func NewQmon(rdb Database) *txqmon {
	//channel for storing mint jobs
	//takes a pointer to a transaction object
	Mintque := make(chan *blockchain.transaction, 3)
	Qmon := &txqmon{
		db:    rdb,
		MintQ: Mintque,
	}
	return Qmon
}

//function to start the transaction que monitoring loop
func (qmon *txqmon) startTXQmonLoop() {
	go qmon.TXloop()
}

//constantly checks the loops
//start jobs by loading into the job queue channel
//
func (qmon *txqmon) TXloop() {
	for {
		//check if the transaction queue has a job
		account, resourceID := qmon.db.DQmint()
		if account != "" {
			//initiate a new transaction
			//load tx into buffered channel
			blockchain.NewTransaction()
		}

	}
}

//loops through storing items into a verfication channel
//the channel should take a unique verfication object for each transaction
func (qmon *txqmon) Verfificationloop() {
	for {
		//check if the transaction queue has a job
		hash, account, resourceID := qmon.db.DQpending()
		if account != "" {
			//initiate a new transaction
			//load tx into buffered channel
			blockchain.NewTransaction()
		}
	}
}

func (qmon *txqmon) txqmanager() {
	//when there are no jobs in the que, kill the manager
	for {
		select {

		case done <- killChan:
			return
		//start a transaction
		case tx <- Minq:

			//start a mint worker
		}

	}
}
