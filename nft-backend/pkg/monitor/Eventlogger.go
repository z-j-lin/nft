package monitor

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogMinted struct {
	from    common.Address
	tokenID *big.Int
}

type LogDeletedTokens struct {
	deleteIds []*big.Int
}

type Events struct{
	
}

func set
