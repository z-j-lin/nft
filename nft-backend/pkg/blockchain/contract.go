package blockchain

import (
	"log"

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
		log.Fatalf("failed to initiate contract instance: %v", err)
	}

	c.instance = instance

}
