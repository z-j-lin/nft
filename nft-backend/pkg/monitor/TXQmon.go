package monitor

import (
	"fmt"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type TXQmon struct {
	db    *redisDb.Database
	eth   *blockchain.Ethereum
	MintQ chan [2]string
	done  chan bool
	QMANI bool
}

func NewQmon(rdb *redisDb.Database, eth *blockchain.Ethereum) *TXQmon {
	//channel for storing mint jobs
	Mintque := make(chan [2]string, 3)
	Qmon := &TXQmon{
		db:    rdb,
		MintQ: Mintque,
		eth:   eth,
	}
	return Qmon
}

//function to start the Transaction que monitoring loop
func (qmon *TXQmon) StartTXQmon() {
	//start the transaction que monitor
	go qmon.TXloop()
}

//checks the mint q for jobs
func (qmon *TXQmon) TXloop() {
	var txinfo [2]string
	numWorkers := make(chan bool, 3)
	for {

		//check if the Transaction queue has a job
		account, resourceID := qmon.db.DQmint()
		if account != "" {
			//numworkers wont be read from until a job is finished
			//blocks when the buffer is full
			numWorkers <- true
			if !qmon.QMANI {
				go qmon.txqmanager(qmon.MintQ, numWorkers)
				qmon.QMANI = true
			}
			//channel the transaction information to the TX manager
			txinfo[0] = account
			txinfo[1] = resourceID
			qmon.MintQ <- txinfo
		}
	}
}

//this works a little too complicated can be simpler, im over it rn
func (qmon *TXQmon) txqmanager(MintQ chan [2]string, taskStatus chan bool) {
	ongoingtasks := 0
	for {
		fmt.Println("task Count:", ongoingtasks)
		select {
		//start a Transaction
		case tx, more := <-MintQ:
			//if nothing is in the mintq kill the manager
			if !more {
				return
			}
			fmt.Println("1", tx[0], tx[1])
			//start a mint worker
			ongoingtasks += 1
			go blockchain.NewTransaction(tx[0], tx[1], qmon.eth.Contract, qmon.db, taskStatus)
		}
	}
}

//loops through storing items into a verfication channel
//the channel should take a unique verfication object for each Transaction
//this should be ran in its own go routine
//function is expected to block when the verification chan is full
/*
func (qmon *TXQmon) Verfificationloop() {
	killChan := make(chan bool)
	for {
		//check if there are pending transactions
		hash, account, resourceID := qmon.db.DQpending()
		if account != "" {
			//check if the verification manager is up
			//if theres a pending transaction
			//pass the job to the verification manager via a buffered channel
		}
	}
}
*/
