package monitor

import (
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type TXQmon struct {
	db       *redisDb.Database
	eth      *blockchain.Ethereum
	MintQ    chan [2]string
	PendingQ chan [3]string
	done     chan bool
	QMANI    bool
}

func NewQmon(rdb *redisDb.Database, eth *blockchain.Ethereum) *TXQmon {
	//channel for storing mint jobs
	Mintque := make(chan [2]string, 10)
	Pendingque := make(chan [3]string, 10)
	Qmon := &TXQmon{
		db:       rdb,
		MintQ:    Mintque,
		PendingQ: Pendingque,
		eth:      eth,
	}
	return Qmon
}

//function to start the Transaction que monitoring loop
func (qmon *TXQmon) StartTXQmon() {
	//start the transaction que monitor
	go qmon.TXloop()
}
func (qmon *TXQmon) QueryMintQ() {
	var txinfo [2]string
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
func (qmon *TXQmon) QueryPendingQ() {
	var txinfo [3]string
	TxHash, recipientAddr, resourceID := qmon.db.DQpending()
	if TxHash != "" {
		if !qmon.QMANI {
			go qmon.PenQManager(qmon.PendingQ)
			qmon.QMANI = true
		}
		//channel the transaction information to the TX manager
		txinfo[0] = account
		txinfo[1] = resourceID
		qmon.MintQ <- txinfo
	}
}

//checks the mint q for jobs
func (qmon *TXQmon) TXloop() {
	//ques := qmon.db.Client.Subscribe(context.TODO())
	//qmessage := ques.Channel()
	for {
		/*
			select{
			case message:= <-qmessage:
				switch message.Payload{
				case "MintQ":
					qmon.QueryMintQ()
				case "PendingQ":
				}
			}
		*/
		//TODO: implement pubsub
		//in the function that uses qmint, publishes a message everytime
		qmon.QueryMintQ()
		//check if there is a pending transaction

		//if there is a pending transaction get the transaction logs
		//if the currrent block is 30 more than the transaction logs
	}
}

//this works a little too complicated can be simpler, im over it rn
func (qmon *TXQmon) txqmanager(MintQ chan [2]string) {
	ongoingtasks := 0
	numWorkers := make(chan bool, 3)
	for {
		//if this channel buffer is full it blocks, no new task is created until a task finishes
		numWorkers <- true
		select {
		//start a Transaction
		case tx, more := <-MintQ:
			//if nothing is in the mintq kill the manager
			if !more {
				qmon.QMANI = false
				return
			}
			//start a mint worker
			ongoingtasks += 1
			go blockchain.NewTransaction(tx[0], tx[1], qmon.eth.Contract, qmon.db, numWorkers)
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
