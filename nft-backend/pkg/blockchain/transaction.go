//This file is Copyright (C) 1997 Master Hacker, ALL RIGHTS RESERVED
package blockchain

import (
	"fmt"
	"log"
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
	resourceID    string
	nonce         *big.Int
	Auth          *bind.TransactOpts
}

//returns the hash of the transaction sent
func NewMintTransaction(TokenRecipient, resourceID string, tnonce *big.Int, eth *Ethereum, rdb *redisDb.Database) Send {
	traddr := common.HexToAddress(TokenRecipient)
	return &mintTx{
		recipientAddr: traddr,
		resourceID:    resourceID,
		eth:           eth,
		nonce:         tnonce,
	}
}

//function to send the transaction
func (mtx *mintTx) SendTransaction(auth *bind.TransactOpts) (*types.Transaction, error) {
	//addr := ethcrypto.PubkeyToAddress(key.PublicKey)
	//sendtx
	log.Printf("Transaction: recipientAddr: %v, resourceID: %s TxNonce: %d\n", mtx.recipientAddr, mtx.resourceID, mtx.nonce)
	tx, err := mtx.eth.Contract.MintToken(auth, mtx.recipientAddr, mtx.resourceID, mtx.nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction:   %v", err)
	} else {
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
