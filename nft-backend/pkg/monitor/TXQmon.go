package monitor

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
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

func NewQmon(redisAddr string, numWorkers int, eth *blockchain.Ethereum, db *redisDb.Database) {
	pm, err := NewPrivKManager()
	if err != nil {
		log.Fatalf("failed to create private key manager: %v", err)
	}
	//add keys to the key manager
	for addr, key := range eth.Keys {
		pm.AddPrivk(addr.Hex(), *key.PrivateKey)
	}
	nm := NewNonceManager(eth)
	if err != nil {
		log.Panic(err)
	}
	//task handler object
	hdl := NewTaskHandler(pm, eth, db, nm)
	NewServerClient(redisAddr, numWorkers, hdl)
}

type NonceMan struct {
	sync.Mutex
	nonce int64
}

func NewNonceManager(eth *blockchain.Ethereum) *NonceMan {
	//get next nonce from contract

	nonce, err := eth.Contract.GetInitNonce()
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

func NewServerClient(redisAddr string, numWorkers int, hdl *Handler) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		//if concurrency is 0, the default would be # accessable CPU
		asynq.Config{Concurrency: 0},
	)
	mux := asynq.NewServeMux()
	// matches the task type with the function to perform the task
	mux.HandleFunc(tasks.TypeMintToken, hdl.HandleMintTokenTask)
	mux.HandleFunc(tasks.TypeBlockVerfication, hdl.HandleVerificationTask)
	//mux.HandleFunc(tasks.TypeBurnTokens, hdl.HandleBurnTokenTask)
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
	nm           *NonceMan
}

//runs as a go routine within the server
//creates a new minttx and txworker object
func NewTaskHandler(PM *PrivkManager, eth *blockchain.Ethereum, db *redisDb.Database, nm *NonceMan) *Handler {
	hdl := &Handler{
		PrivkManager: PM,
		eth:          eth,
		db:           db,
		nm:           nm,
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
	//interface object for sending the mint transaction
	tnonce := hdl.nm.GetnonceWithLock()
	Nonce := big.NewInt(tnonce)
	send := blockchain.NewMintTransaction(data.AccountAddress, data.ResourceID, Nonce, hdl.eth, hdl.db)
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
	//validate the block
	err = NewValidator(hdl.eth, big.NewInt(Data.Blocknum))
	return err
}
func (hdl *Handler) HandleBurnTokenTask(t *asynq.Task) error {
	var Data tasks.BurnToken
	err := json.Unmarshal(t.Payload(), &Data)
	if err != nil {
		log.Println("failed to unmarshal task payload in burnToken task Handler")
		return err
	}
	send := blockchain.NewDelTokens(hdl.eth, Data.TokenIDs)
	NewTX := NewTXWorker(hdl.PrivkManager, hdl.eth, send)
	err = NewTX.Run()
	return err
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

func NewPrivKManager() (*PrivkManager, error) {
	consumedMap := make(map[string]bool)
	availableMap := make(map[string]bool)
	masterSetMap := make(map[string]ecdsa.PrivateKey)
	//New key manager instance
	pm := &PrivkManager{
		consumedMap:  consumedMap,
		availableMap: availableMap,
		masterSetMap: masterSetMap,
	}
	return pm, nil
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
	// make a keyed transactor
	//set transact opts
	auth := txw.init_transactOpt(privk)
	//send the transaction
	tx, err := txw.sendTX.SendTransaction(auth)
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

// given the private key returns a keyed transactor
func (txw *TxWorker) init_transactOpt(privateKey ecdsa.PrivateKey) *bind.TransactOpts {
	pk := &privateKey
	auth, err := bind.NewKeyedTransactorWithChainID(pk, txw.eth.ChainID)
	log.Println("getting auth")
	if err != nil {
		log.Fatal(err)
	}
	// collect the nonce and the gas price
	client := txw.eth.Client
	fromAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//options for transaction
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	return auth
}
