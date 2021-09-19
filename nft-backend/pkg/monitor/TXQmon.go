//This file is Copyright (C) 1997 Master Hacker, ALL RIGHTS RESERVED
package monitor

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	//task handler object
	hdl := NewTaskHandler(pm, eth, db)
	NewServerClient(redisAddr, numWorkers, hdl)
}

func NewServerClient(redisAddr string, numWorkers int, hdl *Handler) {

	TC := tasks.NewTaskClient(hdl.eth, redisAddr)
	TC.QBurnTask()
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		//if concurrency is 0, the default would be # accessable CPU
		asynq.Config{
			Concurrency: 15,
			Queues: map[string]int{
				"transactions": 3,
				"validations":  5,
				"burn":         2,
			},
		},
	)
	mux := asynq.NewServeMux()
	// matches the task type with the function to perform the task
	mux.HandleFunc(tasks.TypeMintToken, hdl.HandleMintTokenTask)
	mux.HandleFunc(tasks.TypeBlockVerfication, hdl.HandleVerificationTask)
	mux.HandleFunc(tasks.TypeBurnTokens, hdl.HandleBurnTokenTask)
	//mux.HandleFunc(tasks.TypeBurnTokens, hdl.HandleBurnTokenTask)
	//starts the server and blocks until a OS signal to exit is sent to terminate
	err := srv.Run(mux)
	if err != nil {
		log.Fatal("unable to start task server", err)
	}
}

//worker with access to private key manager
type TxWorker struct {
	pm     *PrivkManager
	eth    *blockchain.Ethereum
	sendTX blockchain.Send
}

//wrapper object for transactions tasks
//gets the private key from the key manager
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
		log.Println("TXQmon: failed to get privatekey", err)
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
		log.Println("failed to send transaction", err)
		return fmt.Errorf("failed to send transaction %v", err)
	}
	receipt, err := txw.eth.Client.TransactionReceipt(context.TODO(), tx.Hash())
	for receipt == nil && err == ErrTXFailedToRun {
		receipt, err = txw.eth.Client.TransactionReceipt(context.TODO(), tx.Hash())
		//if transaction failed to run on contract
		if err != nil && err != ErrTXFailedToRun {
			log.Println("TXQmon: error @ receipt", err)
			return err
		}
		if receipt.Status == 1 {
			return nil
		}
	}
	return nil
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
