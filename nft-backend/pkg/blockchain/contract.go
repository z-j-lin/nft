//This file is Copyright (C) 1997 Master Hacker, ALL RIGHTS RESERVED
package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
)

type Contract struct {
	ContractAddress common.Address
	eth             *Ethereum
	Instance        *CAToken.CAToken
	MintEvent       string
	TransferEvent   common.Hash
}

func NewContract(eth *Ethereum, ConAddr string) *Contract {
	ContractAddr := common.HexToAddress(ConAddr)
	instance, err := CAToken.NewCAToken(ContractAddr, eth.Client)
	if err != nil {
		log.Fatalf("failed to initiate contract instance: %v", err)
	}
	LogMintedSig := []byte("Minted(address,uint256)")
	LogTokenTransferSig := []byte("Transfer(address,address,uint256)")
	LogMintedSigHash := crypto.Keccak256Hash(LogMintedSig).Hex()
	LogTokenTransferSigHash := crypto.Keccak256Hash(LogTokenTransferSig)
	con := &Contract{
		MintEvent:       LogMintedSigHash,
		TransferEvent:   LogTokenTransferSigHash,
		Instance:        instance,
		eth:             eth,
		ContractAddress: ContractAddr,
	}
	return con
}
func (c *Contract) MintToken(Auth *bind.TransactOpts, RecipientAddr common.Address, resourceID string, nonce *big.Int) (*types.Transaction, error) {
	tx, err := c.Instance.Mint(Auth, RecipientAddr, resourceID, nonce)
	if err != nil {
		log.Println("failed to send transaction @ MintToken:", err)
		return nil, err
	}
	return tx, err
}

//gets the next unused nonce from the contract
//only use for inital startup
//this is not concurrency safe
func (c *Contract) GetInitNonce() (*big.Int, error) {
	opts := &bind.CallOpts{
		Pending: false,
		Context: context.TODO(),
	}
	nonce, err := c.Instance.CATokenCaller.NextNonce(opts)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

func (c *Contract) DeleteTokens(auth *bind.TransactOpts, IDs []*big.Int) (*types.Transaction, error) {
	tx, err := c.Instance.ExpiredContracts(auth, IDs)
	if err != nil {
		log.Println("failed to send transaction @ DeleteTokens:", err)
		return nil, err
	}
	return tx, err
}

func (c *Contract) SetServerRole(Auth *bind.TransactOpts, serverAddress common.Address) (*types.Transaction, error) {
	tx, err := c.Instance.AddServerRole(Auth, serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	return tx, err
}

//check if owner owns token
func (c *Contract) IsOwner(address string, tokenID string) (bool, error) {
	//TODO: use highest delayed block by 20 blocks
	opts := &bind.CallOpts{
		Pending: false,
		Context: context.TODO(),
	}

	TokenID, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return false, err
	}
	TID := big.NewInt(TokenID)
	ownerAddr, err := c.Instance.OwnerOf(opts, TID)
	if err != nil {
		return false, err
	}
	account := common.HexToAddress(address)
	if ownerAddr == account {
		return true, nil
	} else {
		return false, fmt.Errorf("not owner")
	}
}

func (c *Contract) GetResourceID(tokenID string) (string, error) {
	//TODO: use highest delayed block by 20 blocks
	opts := &bind.CallOpts{
		Pending: false,
		Context: context.TODO(),
	}
	TokenID, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return "", err
	}
	tokenIDBI := big.NewInt(TokenID)
	TokenURI, err := c.Instance.TokenURI(opts, tokenIDBI)
	return TokenURI, err
}

func (c *Contract) SafeTransferFrom(Auth *bind.TransactOpts, Owneraddress, Recipient common.Address, tokenID int64) (*types.Transaction, error) {

	TID := big.NewInt(tokenID)
	tx, err := c.Instance.SafeTransferFrom(Auth, Owneraddress, Recipient, TID)
	return tx, err
}
