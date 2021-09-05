package monitor

import (
	"math/big"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type TXQmon struct {
	db            *redisDb.Database
	eth           *blockchain.Ethereum
	MintQ         chan [2]string
	PendingBlockQ chan uint64
	QMANI         bool
	QBMANI        bool
}

func NewQmon(rdb *redisDb.Database, eth *blockchain.Ethereum) *TXQmon {
	//channel for storing mint jobs
	Mintque := make(chan [2]string, 10)
	PendingBlockque := make(chan uint64, 10)
	Qmon := &TXQmon{
		db:            rdb,
		MintQ:         Mintque,
		PendingBlockQ: PendingBlockque,
		eth:           eth,
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
		qmon.QueryPendingBlockQ()
		//if there is a pending transaction get the transaction logs
	}
}

//function to query the mintQ
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

//this works a little too complicated can be simpler, im over it rn
func (qmon *TXQmon) txqmanager(MintQ chan [2]string) {
	numWorkers := make(chan bool, 3)
	for {
		//if this channel buffer is full it blocks, no new task is created until a task finishes
		numWorkers <- true
		select {
		//start a Transaction
		case tx := <-MintQ:
			//if nothing is in the mintq kill the manager
			//start a mint worker
			go blockchain.NewTransaction(tx[0], tx[1], qmon.eth.Contract, qmon.db, numWorkers)
		default:
			qmon.QMANI = false
			return
		}
	}
}

//query the block Validator Q
func (qmon *TXQmon) QueryPendingBlockQ() {
	Blocknum := qmon.db.DQpendingBlock()
	if Blocknum != 0 {
		if !qmon.QBMANI {
			go qmon.ValManager(qmon.PendingBlockQ)
			qmon.QBMANI = true
		}
		//add blocknum to the buffer
		qmon.PendingBlockQ <- Blocknum
	}
}
func (qmon *TXQmon) ValManager(BlockQ chan uint64) {
	numWorkers := make(chan bool, 8)
	for {
		//if this channel buffer is full it blocks, no new task is created until a task finishes
		numWorkers <- true
		select {
		//start a Transaction
		case Blocknum := <-BlockQ:
			//start a mint worker
			num := big.NewInt(int64(Blocknum))
			go NewValidator(qmon.eth, num, numWorkers)
		default:
			//if nothing is in the mintq kill the manager
			qmon.QBMANI = false
			return
		}
	}
}
