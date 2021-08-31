package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
)

type MintTx struct {
	contract      *Contract
	Auth          *bind.TransactOpts
	recipientAddr common.Address
	db            *redisDb.Database
	resourceID    string
}

//returns a pointer to a new transaction object
func NewTransaction(TokenRecipient, resourceID string, contract *Contract, rdb *redisDb.Database, taskStatus chan bool) {
	//instantiate new keyed transactor
	auth := bind.NewKeyedTransactor(contract.eth.Key.PrivateKey)
	traddr := common.HexToAddress(TokenRecipient)
	tranx := &MintTx{
		Auth:          auth,
		contract:      contract,
		recipientAddr: traddr,
		db:            rdb,
	}
	tranx.resourceID = resourceID
	fmt.Println("sending")

	tranx.SendTransaction(TokenRecipient, taskStatus)
	return
}

func (mtx *MintTx) init_transactOpt() {
	// collect the nonce and the gas price
	auth := mtx.Auth
	client := mtx.contract.eth.Client
	fromAddress := mtx.contract.eth.Account.Address
	fmt.Println(fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	//options for transaction
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
}

//function to send the transaction
func (mtx *MintTx) SendTransaction(address string, taskStatus chan bool) {

	to := common.HexToAddress(address)
	Catoken := mtx.contract.instance
	mtx.init_transactOpt()
	//set true for testing monitor
	mtx.Auth.NoSend = true
	//sendtx
	tx, err := Catoken.Mint(mtx.Auth, to)
	fmt.Println("qpending")
	if err != nil {
		log.Printf("transaction failed: %v", err)
		//add the transaction back to the transaction que
		mtx.db.Qmint(mtx.recipientAddr.Hex(), mtx.resourceID)
	} else {
		//if didnt fail add the transaction to the pending list
		fmt.Println("qpending")
		mtx.db.Qpending(tx.Hash().Hex(), address, mtx.resourceID)
		taskStatus <- true
		fmt.Println("after qpending")
	}
}
