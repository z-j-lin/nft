package monitor

import (
	"math/big"
	"testing"

	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

func TestReadevents(t *testing.T) {

	//instantiate ether instance
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	contractAddr := "0xb410756d52b1250aB9bE358437Ab41a4D7636Af8"
	eth, err := blockchain.NewEtherClient(rpcurl, contractAddr, big.NewInt(int64(3)))
	if err != nil {
		t.Error(err)
	}
	blocknum := big.NewInt(int64(11049890))
	NewValidator(eth, blocknum)
}
