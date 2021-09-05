package blockchain

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestDeleteTokens(t *testing.T) {
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	contractAddr := "0x3DA85558aF6D0d0D03283fa23eD1edE90f7E3E03"
	eth, err := NewEtherClient(rpcurl, contractAddr)
	if err != nil {
		panic(err)
	}
	con := eth.Contract
	auth := bind.NewKeyedTransactor(eth.Key.PrivateKey)
	var tokens []*big.Int
	tokens = append(tokens, big.NewInt(int64(0)), big.NewInt(int64(1)))
	con.DeleteTokens(auth, tokens)
}
