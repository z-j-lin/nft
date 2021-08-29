package blockchain

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type transaction struct {
	contract      *Contract
	Auth          *bind.TransactOpts
	recipientAddr common.Address
}

func NewTransaction(TokenRecipient string, contract *Contract) *transaction {
	//instantiate new keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	traddr := common.HexToAddress(TokenRecipient)
	return &transaction{
		Auth:          auth,
		contract:      contract,
		recipientAddr: traddr,
	}

}

func (tx *transaction) init_transactOpt() {
	// collect the nonce and the gas price
	auth := tx.Auth
	client := tx.contract.eth.client
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background(), fromAddress)
	//options for transaction

}

//
func (tx *transaction) SendTransactions() {
	Catoken := tx.contract.instance
	tx.init_transactOpt()
	Catoken.Mint(tx.Auth)

}
