package monitor

import (
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type TXQmon struct {
	db    *redisDb.Database
	eth   *blockchain.Ethereum
	MintQ chan []string
	QMANI bool
}

func NewQmon(rdb *redisDb.Database) *TXQmon {
	//channel for storing mint jobs
	//takes a pointer to a Transaction object
	Mintque := make(chan [2]string, 3)

	Qmon := &TXQmon{
		db:    rdb,
		MintQ: Mintque,
	}
	return Qmon
}

//function to start the Transaction que monitoring loop
func (qmon *TXQmon) startTXQmon() {
	go qmon.TXloop()
}

//checks the mint q for jobs
//blocks when MintQ chan is full
//start jobs by loading into the job queue channel
//
func (qmon *TXQmon) TXloop() {
	var txinfo [2]string
	for {
		//check if the Transaction queue has a job
		account, resourceID := qmon.db.DQmint()
		if account != "" {
			if !qmon.QMANI {
				go qmon.txqmanager()
			}
			//channel the transaction information to the TX manager
			txinfo[0] = account
			txinfo[1] = resourceID
			qmon.MintQ <- txinfo[:]
		}

	}
}

//loops through storing items into a verfication channel
//the channel should take a unique verfication object for each Transaction
func (qmon *TXQmon) Verfificationloop() {
	for {
		//check if the Transaction queue has a job
		hash, account, resourceID := qmon.db.DQpending()
		if account != "" {

		}
	}
}

func (qmon *TXQmon) txqmanager() {
	//when there are no jobs in the que, kill the manager
	done := false
	for !done {
		select {

		case done <- killChan:
			return
		//start a Transaction
		case tx <- Mintq:

			//start a mint worker
		}

	}
}
