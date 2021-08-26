package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	CAToken "./build/CAToken"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	client, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Panic("cant connect to rpc server: ", err)
	}
	//test deployment account
	privateKey, err := crypto.HexToECDSA("f43975530a55cd2f17d1973750b4b803fa651067a222de74a16f60864d0c452e")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")

	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	serverAddress := common.HexToAddress("0x282e0869A1fA76B7f6d2D4B2758b60bF0284F418")
	address, tx, instance, err := CAToken.DeployCAToken(auth, client, "CAToken", "CAT", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())
	_ = instance
}
