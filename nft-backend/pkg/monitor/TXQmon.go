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
				go qmon.txqmanager(qmon.MintQ)
				qmon.QMANI = true
			}
			//channel the transaction information to the TX manager
			txinfo[0] = account
			txinfo[1] = resourceID
			qmon.MintQ <- txinfo
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
func (qmon *TXQmon) txqmanager(MintQ chan [2]string) {
	taskStatus := make(chan bool, 2)
	ongoingtasks := 0

	for {
		fmt.Println("task Count:", ongoingtasks)
		if ongoingtasks < 3 {
			select {
			//start a Transaction
			case tx := <-MintQ:
				fmt.Println("1", tx[0], tx[1])
				//start a mint worker
				ongoingtasks += 1
				go blockchain.NewTransaction(tx[0], tx[1], qmon.eth.Contract, qmon.db, taskStatus)
			}
		}
		if ongoingtasks > 2 {
			select {
			//if nothing is in the mintq kill the manager
			case <-taskStatus:
				fmt.Println("ongoingtask at decrementer", ongoingtasks)
				ongoingtasks -= 1
			}
		}
	}
}
