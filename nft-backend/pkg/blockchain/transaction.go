package blockchain

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
)

type Transaction struct {
	contract      *Contract
	Auth          *bind.TransactOpts
	recipientAddr common.Address
	db            *redisDb.Database
	resourceID    string
}

//returns a pointer to a new transaction object
func NewTransaction(TokenRecipient, resourceID string, contract *Contract) *Transaction {
	//instantiate new keyed transactor
	auth := bind.NewKeyedTransactor(contract.eth.key.PrivateKey)
	traddr := common.HexToAddress(TokenRecipient)
	return &Transaction{
		Auth:          auth,
		contract:      contract,
		recipientAddr: traddr,
	}
}

func (tx *Transaction) init_transactOpt() {
	// collect the nonce and the gas price
	auth := tx.Auth
	client := tx.contract.eth.client
	fromAddress := tx.contract.eth.account.Address
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
func (tx *Transaction) SendTransaction(address, resourceID string) {
	to := common.HexToAddress(address)
	Catoken := tx.contract.instance
	tx.init_transactOpt()
	receipt, err := Catoken.Mint(tx.Auth, to)
	if err != nil {
		log.Printf("transaction failed: %v", err)
	}
	//wait 10 seconds, check if its through
	//if failed exit
	txhash := receipt.Hash().Hex()
	//if didnt fail add the transaction to the pending list
	tx.db.Qpending(txhash, address, resourceID)
}
