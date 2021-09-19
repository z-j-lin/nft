package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/hibiken/asynq"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

var ErrNoKeys error = errors.New("no privk available")
var ErrKeyConflict error = errors.New("privk added twice")
var ErrTXFailedToRun error = errors.New("not found")

type Handler struct {
	PrivkManager *PrivkManager
	eth          *blockchain.Ethereum
	db           *redisDb.Database
	nm           *NonceMan
	TC           *tasks.TaskClient
}

//runs as a go routine within the server
//creates a new minttx and txworker object
func NewTaskHandler(PM *PrivkManager, eth *blockchain.Ethereum, db *redisDb.Database) *Handler {
	tc := tasks.NewTaskClient(db.Client.Options().Addr)
	nm := NewNonceManager(eth)
	hdl := &Handler{
		PrivkManager: PM,
		eth:          eth,
		db:           db,
		nm:           nm,
		TC:           tc,
	}
	return hdl
}
func (hdl *Handler) HandleMintTokenTask(ctx context.Context, t *asynq.Task) error {
	//data struct stores data for the task
	var data tasks.MintToken
	err := json.Unmarshal(t.Payload(), &data)
	if err != nil {
		log.Println("failed to unmarshal task payload in Minttoken Handler")
		return err
	}
	//get mint transaction nonce
	tnonce := hdl.nm.GetnonceWithLock()
	//interface object for sending the mint transaction
	send := blockchain.NewMintTransaction(data.AccountAddress, data.ResourceID, big.NewInt(tnonce), hdl.eth, hdl.db)
	//start a new txworker, with the minttx object
	NewTX := NewTXWorker(hdl.PrivkManager, hdl.eth, send)
	//run the worker
	err = NewTX.Run()
	//returns the status of the job
	return err
}
func (hdl *Handler) HandleVerificationTask(ctx context.Context, t *asynq.Task) error {
	var Data tasks.BlockV
	err := json.Unmarshal(t.Payload(), &Data)
	if err != nil {
		log.Println("failed to unmarshal task payload in Verification Handler")
		return err
	}
	//instantiates a new validator in a go routine started by the asynq server
	err = NewValidator(hdl.eth, big.NewInt(Data.Blocknum))
	return err
}

func (hdl *Handler) HandleBurnTokenTask(t *asynq.Task) error {
	//gets the expired tokens from the db
	tokens, err := hdl.db.GetExpiredTokens()
	//container for tokenIDs
	var Tokens []*big.Int
	//convert string array to big int array
	//theres got to be a better way to do this
	for _, t := range tokens {
		//convert each string into a integeger
		n := new(big.Int)
		n, ok := n.SetString(t, 10)
		if !ok {
			return fmt.Errorf("failed to set token string to bigint")
		}
		Tokens = append(Tokens, n)
	}
	//run the tx if there are tokens in the array
	if len(Tokens) != 0 {
		//make the transaction object
		send := blockchain.NewDelTokens(hdl.eth, Tokens)
		//spawn a tx worker with the transaction object
		NewTX := NewTXWorker(hdl.PrivkManager, hdl.eth, send)
		//run the transaction
		err := NewTX.Run()
		//delete tokens from the db record after it gets deleted
		if err == nil {
			err = hdl.db.DeleteExpiredTokens()
		}

	}
	//make another task to run later
	err = hdl.TC.QBurnTask()
	return err
}

func NewNonceManager(eth *blockchain.Ethereum) *NonceMan {
	//get next nonce from contract
	nonce, err := eth.Contract.GetInitNonce()
	log.Println("TXQmon: started new nonce manager with nonce", nonce)
	if err != nil {
		log.Panic(err)
	}
	return &NonceMan{
		nonce: nonce.Int64(),
	}
}
func (nm *NonceMan) GetnonceWithLock() int64 {
	nm.Lock()
	defer nm.Unlock()
	nonce := nm.nonce
	nm.nonce = nonce + 1
	return nonce
}
