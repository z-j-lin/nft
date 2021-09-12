package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
)

type Send interface {
	SendTransaction(*bind.TransactOpts) (*types.Transaction, error)
}
type mintTx struct {
	eth           *Ethereum
	recipientAddr common.Address
	db            *redisDb.Database
	resourceID    string
	nonce         *big.Int
	Auth          *bind.TransactOpts
}
type burnTx struct {
}

//returns the hash of the transaction sent
func NewMintTransaction(TokenRecipient, resourceID string, tnonce *big.Int, eth *Ethereum, rdb *redisDb.Database) Send {
	traddr := common.HexToAddress(TokenRecipient)
	return &mintTx{
		recipientAddr: traddr,
		resourceID:    resourceID,
		eth:           eth,
		db:            rdb,
		nonce:         tnonce,
	}
}

//function to send the transaction
func (mtx *mintTx) SendTransaction(auth *bind.TransactOpts) (*types.Transaction, error) {
	//addr := ethcrypto.PubkeyToAddress(key.PublicKey)
	//sendtx
	tx, err := mtx.eth.Contract.MintToken(auth, mtx.recipientAddr, mtx.nonce)
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

type DelTokens struct {
	Tokens []*big.Int
	eth    *Ethereum
}

func NewDelTokens(eth *Ethereum, Toks []*big.Int) Send {
	return &DelTokens{
		Tokens: Toks,
		eth:    eth,
	}
}
func (dt *DelTokens) SendTransaction(auth *bind.TransactOpts) (*types.Transaction, error) {
	tx, err := dt.eth.Contract.DeleteTokens(auth, dt.Tokens)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
