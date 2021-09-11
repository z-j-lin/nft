package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
)

type Send interface {
	SendTransaction(ecdsa.PrivateKey) (*types.Transaction, error)
}
type mintTx struct {
	eth           *Ethereum
	recipientAddr common.Address
	db            *redisDb.Database
	resourceID    string
}
type burnTx struct {
}

//returns the hash of the transaction sent
func NewMintTransaction(TokenRecipient, resourceID string, eth *Ethereum, rdb *redisDb.Database) Send {
	traddr := common.HexToAddress(TokenRecipient)
	tranx := &mintTx{
		recipientAddr: traddr,
		resourceID:    resourceID,
		eth:           eth,
		db:            rdb,
	}
	return tranx
}

func NewBurnTransaction() {

}

func (mtx *mintTx) init_transactOpt(privateKey ecdsa.PrivateKey) *bind.TransactOpts {
	pk := &privateKey
	auth, err := bind.NewKeyedTransactorWithChainID(pk, mtx.eth.chainID)
	log.Println("getting auth")
	if err != nil {
		log.Fatal(err)
	}
	// collect the nonce and the gas price
	client := mtx.eth.Client
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

//function to send the transaction
func (mtx *mintTx) SendTransaction(key ecdsa.PrivateKey) (*types.Transaction, error) {
	//addr := ethcrypto.PubkeyToAddress(key.PublicKey)
	auth := mtx.init_transactOpt(key)
	//set true for testing monitor
	auth.NoSend = false
	//sendtx
	tx, err := mtx.eth.Contract.MintToken(auth, mtx.recipientAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction:   %v", err)
	} else {
		//if didnt fail add the transaction to the pending list
		//temp map of resource ID to txn hash
		mtx.db.Client.Set(context.TODO(), tx.Hash().Hex(), mtx.resourceID, 0)
		fmt.Println("transaction added to pending que")
		return tx, nil
	}

}
