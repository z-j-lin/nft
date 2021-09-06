package blockchain

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func TestDeleteTokens(t *testing.T) {

	keybyte := common.Hex2Bytes("23d40b2eeb9bd399a8a310c67f5fc5233c1ff734f4a3c68f0cbd488836b64d01")
	key, err := ethcrypto.ToECDSA(keybyte)
	if err != nil {
		t.Fatal(err)
	}
	senderAddr := ethcrypto.PubkeyToAddress(key.PublicKey)
	t.Logf("%x", senderAddr)
	chainID := big.NewInt(int64(5777))
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		t.Fatal(err)
	}
	rpcurl := "http://127.0.0.1:9545"
	contractAddr := "0xE73bF47beA2FA0d70Dbaf0316C38E9f42881f940"
	eth, err := NewEtherClient(rpcurl, contractAddr)
	if err != nil {
		t.Fatal(err)
	}
	con := eth.Contract

	//fromAddress := common.HexToAddress("0x7B725F2ae9e159ADD3D49DbEA88C841e8fC52793")
	nonce, err := eth.Client.PendingNonceAt(context.Background(), senderAddr)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	//options for transaction
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	var tokens []*big.Int
	tokens = append(tokens, big.NewInt(int64(0)))
	tx0, err := con.MintToken(auth, senderAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", tx0.Hash())

	nonce, err = eth.Client.PendingNonceAt(context.Background(), senderAddr)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	tx1, err := con.DeleteTokens(auth, tokens)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", tx1.Hash())
}
