package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
)

type Contract struct {
	contractAddress common.Address
	eth             *Ethereum
	instance        *CAToken.CAToken
}

func (c *Contract) init() {
	instance, err := CAToken.NewCAToken(c.contractAddress, c.eth.Client)
	if err != nil {
		panic(err)
	}
	c.instance = instance

}
