package blockchain

import (
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewdeployEtherClient(t *testing.T) {
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"

	eth, err := NewdeployEtherClient(rpcurl, big.NewInt(int64(3)))
	if err != nil {
		log.Panic("cant connect to rpc server: ", err)
	}

	for addr, key := range eth.Keys {
		pk := crypto.FromECDSA(key.PrivateKey)
		//keystring := string(pk)
		fmt.Printf("address: %v: %x \n", addr, pk)
	}
}
