package eth

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	CAToken "../CAToken"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Contract struct {
	client          *ethclient.Client
	rpcurl          string
	contractAddress string
}

func (c *Contract) LoadContract(string) {
	client, err := ethclient.Dial(c.rpcurl)
	if err != nil {
		log.Fatalf("Failed to connect to the ether network: %v", err)
	}
	c.client = client
	address := common.HexToAddress(c.contractAddress)
	ct, err := CAToken.NewCAToken(address, c.client)
	if err != nil {
		log.Fatalf("error loading contract %v", err)
	}
}

//prepares a mint transaction
func (c *Contract) BuildTransactionMint(pk string) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatalf("private key Hex2ECDSA error %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("error getting nonce %v", err)
	}
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("error on suggest gas price %v", err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

}
