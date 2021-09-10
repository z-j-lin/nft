package blockchain

import (
	"fmt"
	"log"
	"math/big"

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
	DeleteEvent     string
}

func (c *Contract) init() {
	instance, err := CAToken.NewCAToken(c.ContractAddress, c.eth.Client)
	if err != nil {
		log.Fatalf("failed to initiate contract instance: %v", err)
	}
	LogMintedSig := []byte("Minted(address,uint256)")
	LogDeletedTokensSig := []byte("DeletedTokens(uint256[])")
	LogMintedSigHash := crypto.Keccak256Hash(LogMintedSig).Hex()
	LogDeletedTokensSigHash := crypto.Keccak256Hash(LogDeletedTokensSig).Hex()
	c.MintEvent = LogMintedSigHash
	c.DeleteEvent = LogDeletedTokensSigHash
	c.Instance = instance

}
func (c *Contract) MintToken(Auth *bind.TransactOpts, RecipientAddr common.Address) (*types.Transaction, error) {
	tx, err := c.Instance.Mint(Auth, RecipientAddr)
	if err != nil {
		log.Println("failed to send transaction @ MintToken:", err)
		return nil, err
	}
	return tx, err
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
func (c *Contract) IsOwner(CallOpts *bind.CallOpts, address common.Address, tokenID int64) (bool, error) {
	TID := big.NewInt(tokenID)
	ownerAddr, err := c.Instance.OwnerOf(CallOpts, TID)
	if err != nil {
		return false, err
	}
	if ownerAddr == address {
		return true, nil
	} else {
		return false, fmt.Errorf("not owner")
	}
}

func (c *Contract) SafeTransferFrom(Auth *bind.TransactOpts, Owneraddress, Recipient common.Address, tokenID int64) (*types.Transaction, error) {
	TID := big.NewInt(tokenID)
	tx, err := c.Instance.SafeTransferFrom(Auth, Owneraddress, Recipient, TID)
	return tx, err
}
