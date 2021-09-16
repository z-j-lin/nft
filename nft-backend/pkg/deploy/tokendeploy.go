package deploy

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

func Deploy() {
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"

	eth, err := blockchain.NewdeployEtherClient(rpcurl, big.NewInt(int64(3)))
	if err != nil {
		log.Panic("cant connect to rpc server: ", err)
	}
	//test deployment account
	ownerAddr := common.HexToAddress("0x359aa05C01338C83A5835BEbC1E689e129a06868")
	privateKey := eth.Keys[ownerAddr].PrivateKey
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")

	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := eth.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, eth.ChainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(5574643) // in units
	auth.GasPrice = gasPrice
	auth.NoSend = false
	serverAddress := common.HexToAddress("0x02a584a34b78b6645eb9728cfbf0e56433e3585b")
	address, tx, instance, err := CAToken.DeployCAToken(auth, eth.Client, "CAToken", "CAT", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())

	fmt.Println(tx.Hash().Hex())
	fmt.Println(tx.Gas())
	_ = instance
}
