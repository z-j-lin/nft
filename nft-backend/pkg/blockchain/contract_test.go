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
	keybyte := common.Hex2Bytes("fabb40d41c3044e3d1200ebb726192d5dde3e349565e27bb6f900556cfacbbe5")
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
	contractAddr := "0x810dA0c61C3b19087d40cdCa990790351F146dc8"
	eth, err := NewEtherClient(rpcurl, contractAddr, chainID)
	if err != nil {
		t.Fatal(err)
	}
	con := eth.Contract
	//fromAddress := common.HexToAddress("0x7B725F2ae9e159ADD3D49DbEA88C841e8fC52793")
	gasPrice, err := eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	//options for transaction

	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	var tokens []*big.Int
	tokens = append(tokens, big.NewInt(int64(0)))
	for i := 0; i < 3; i++ {
		nonce, err := eth.Client.PendingNonceAt(context.Background(), senderAddr)
		if err != nil {
			t.Fatal(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))
		tx0, err := con.MintToken(auth, senderAddr)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%x", tx0.Hash())
	}

	nonce, err := eth.Client.PendingNonceAt(context.Background(), senderAddr)
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

func TestSetServerRole(t *testing.T) {
	rpcurl := "http://127.0.0.1:9545"
	chainID := big.NewInt(int64(5777))
	contractAddr := "0x810dA0c61C3b19087d40cdCa990790351F146dc8"
	eth, err := NewEtherClient(rpcurl, contractAddr, chainID)
	if err != nil {
		t.Fatal(err)
	}
	ownerAcc := common.HexToAddress("0x3BABA5d61c49a21e7C54E25D6b6be4fd1DaA1D7E")
	key := eth.Keys[ownerAcc]
	ownerpk := key.PrivateKey
	senderAddr := ethcrypto.PubkeyToAddress(ownerpk.PublicKey)
	auth, err := bind.NewKeyedTransactorWithChainID(ownerpk, chainID)
	if err != nil {
		log.Panic(err)
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)

	accounts := eth.Accounts
	con := eth.Contract
	for _, account := range accounts {
		gasPrice, err := eth.Client.SuggestGasPrice(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		auth.GasPrice = gasPrice
		nonce, err := eth.Client.PendingNonceAt(context.Background(), senderAddr)
		if err != nil {
			t.Fatal(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))
		con.SetServerRole(auth, account.Address)
	}
}
