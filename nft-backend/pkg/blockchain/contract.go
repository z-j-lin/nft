package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/z-j-lin/nft/blob/main/nft-backend/pkg/monitor"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
)

type Contract struct {
	contractAddress common.Address
	eth             *ethereum
	instance        *CAToken.CAToken
}
