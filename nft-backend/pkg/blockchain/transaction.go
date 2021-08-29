package blockchain

import (
	"context"
	"log"
	"math/big"

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
	bind.NewTransactor()
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
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
}

//
func (tx *transaction) SendTransactions() {
	Catoken := tx.contract.instance
	tx.init_transactOpt()
	Catoken.Mint(tx.Auth)

}
