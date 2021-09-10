package monitor

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hibiken/asynq"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

/*

user -> tx send

			ATOMIC
user -> [srvr -> db as pending q] -> [Qmon -> tx send -> block until mined or failure]

blockMon -> addsToDBForAccess

user w/ proof -> srvr -> validates against chain -> returns data

*/

func NewQmon(redisAddr string, numWorkers int, eth *blockchain.Ethereum) {
	consumedMap := make(map[string]bool)
	availableMap := make(map[string]bool)
	masterSetMap := make(map[string]ecdsa.PrivateKey)
	//New key manager instance
	pm := &PrivkManager{
		consumedMap:  consumedMap,
		availableMap: availableMap,
		masterSetMap: masterSetMap,
	}
	//add keys to the key manager
	for addr, key := range eth.Keys {
		pm.AddPrivk(addr.Hex(), *key.PrivateKey)
	}
	//used in the server client handler functions
	hdl := &Handler{
		PrivkManager: pm,
	}
	NewServerClient(redisAddr, numWorkers, hdl)

}

func NewServerClient(redisAddr string, numWorkers int, hdl *Handler) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		//if concurrency is 0, the default would be # accessable CPU
		asynq.Config{Concurrency: 0},
	)
	mux := asynq.NewServeMux()
	// matches the task type with the function to perform the task
	mux.HandleFunc(tasks.TypeMintToken, hdl.HandleMintTokenTask)
	//mux.HandleFunc(TypeBurnTokens)
	//starts the server and blocks until a OS signal to exit is sent to terminate
	err := srv.Run(mux)
	if err != nil {
		log.Fatal("unable to start task server", err)
	}
}

type Handler struct {
	PrivkManager *PrivkManager
	eth          *blockchain.Ethereum
	db           *redisDb.Database
}

//runs as a go routine within the server
func (hdl *Handler) HandleMintTokenTask(ctx context.Context, t *asynq.Task) error {
	//data struct stores data for the task
	var data tasks.MintToken
	err := json.Unmarshal(t.Payload(), &data)
	if err != nil {
		log.Println("failed to unmarshal task payload in Minttoken Handler")
		return err
	}
	send := blockchain.NewMintTransaction(data.AccountAddress, data.ResourceID, hdl.eth, hdl.db)
	//start a new txworker, with the mint transaction function
	NewTX := TxWorker{
		pm:     hdl.PrivkManager,
		eth:    hdl.eth,
		sendTX: send,
	}
	//run the worker
	err = NewTX.Run()
	//returns the status of the job
	return err
}
func (hdl *Handler) HandleBurnTokenTask(t *asynq.Task) error {
	panic("unimplemented")
	return nil
}

var ErrNoKeys error = errors.New("no privk available")
var ErrKeyConflict error = errors.New("privk added twice")
var ErrTXFailedToRun error = errors.New("not found")

// PrivkManager releases keys to
type PrivkManager struct {
	sync.Mutex
	// stores all possibe private keys OR knows where to go get them
	// string key is the ether account
	consumedMap  map[string]bool
	availableMap map[string]bool
	masterSetMap map[string]ecdsa.PrivateKey
}

func (pm *PrivkManager) AddPrivk(addr string, privk ecdsa.PrivateKey) error {
	pm.Lock()
	defer pm.Unlock()
	_, ok := pm.masterSetMap[addr]
	if ok {
		return ErrKeyConflict
	}
	pm.masterSetMap[addr] = privk
	pm.availableMap[addr] = true
	return nil
}

func (pm *PrivkManager) GetWithLock() (ecdsa.PrivateKey, func(), error) {
	pm.Lock()
	defer pm.Unlock()
	var privkAddr string
	for _privkAddr := range pm.availableMap {
		privkAddr = _privkAddr
		break
	}
	if privkAddr == "" {
		return ecdsa.PrivateKey{}, nil, ErrNoKeys
	}

	privk := pm.masterSetMap[privkAddr]
	delete(pm.availableMap, privkAddr)
	pm.consumedMap[privkAddr] = true
	return privk, pm.free(privkAddr), nil
}

func (pm *PrivkManager) free(privkAddr string) func() {
	return func() {
		//does it need mutex lock? each worker will have a unique key
		pm.Lock()
		defer pm.Unlock()
		delete(pm.consumedMap, privkAddr)
		pm.availableMap[privkAddr] = true
	}
}

//worker with access to private key manager
type TxWorker struct {
	pm     *PrivkManager
	eth    *blockchain.Ethereum
	sendTX blockchain.Send
}

func NewTXWorker(privKM *PrivkManager, eth *blockchain.Ethereum, send blockchain.Send) *TxWorker {
	return &TxWorker{
		pm:     privKM,
		eth:    eth,
		sendTX: send,
	}
}
func (txw *TxWorker) Run() error {
	//get the private key
	privk, free, err := txw.pm.GetWithLock()
	if err != nil {
		return err
	}
	//the key is freed after the transaction is mined
	defer free()
	//send the transaction
	tx, err := txw.sendTX.SendTransaction(privk)
	// if failed to send transaction
	if err != nil {
		//return an error, key is freed, task will retry
		return fmt.Errorf("failed to send transaction")
	}
	receipt, err := txw.eth.Client.TransactionReceipt(context.TODO(), tx.Hash())
	rx := make(chan *types.Receipt)
	for receipt == nil && err == ErrTXFailedToRun {
		receipt, err = txw.eth.Client.TransactionReceipt(context.TODO(), tx.Hash())
		//if transaction failed to run on contract
		if err != nil && err != ErrTXFailedToRun {

			return err
		}
		//load the receipt into the channel
		rx <- receipt
		select {
		case Receipt := <-rx:
			if Receipt.Status == 1 {
				return nil
			} else {
				return fmt.Errorf("transaction failed")
			}
		case <-time.After(10 * time.Minute):
			return fmt.Errorf("timedout failed to send transaction")
		}
	}
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}
