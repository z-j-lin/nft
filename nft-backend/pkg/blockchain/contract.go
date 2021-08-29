package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
)

type Contract struct {
	contractAddress common.Address
	eth             *ethereum
	instance        *CAToken.CAToken
}

func (c *Contract) init() {
	instance, err := CAToken.NewCAToken(c.contractAddress, c.eth.client)
	if err != nil {
		panic(err)
	}
	c.instance = instance

}

func (c *Contract) MintNewtoken(
	Auth *bind.TransactOpts,
	recipientAddress common.Address) (*types.Transaction, error) {
	tx, err := c.instance.Mint(Auth, recipientAddress)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
