package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	CAToken "../CAToken"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Contract struct {
	client  *ethclient.Client
	ConAddr string
	CAToken *CAToken.CAToken
}

func initContract(rpcurl, contractAddress string) (*Contract, error) {
	conn, err := ethclient.Dial(rpcurl)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the ether network: %v", err)
	}

	address := common.HexToAddress(contractAddress)
	CAToken, err := CAToken.NewCAToken(address, conn)
	if err != nil {
		return nil, fmt.Errorf("error loading contract %v", err)
	}
	return &Contract{
		client:  conn,
		ConAddr: contractAddress,
		CAToken: CAToken,
	}, nil
}

type ServerAcc struct {
	auth          *bind.TransactOpts
	WalletAddress common.Address
}

func NewAcc() *ServerAcc {
	//import private key here
	privateKey, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatalf("private key Hex2ECDSA error %v", err)
	}
	publicKey := privateKey.Public()
	//store this in a struct
	authentication := bind.NewKeyedTransactor(privateKey)
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	serverWalletAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &ServerAcc{
		auth:          authentication,
		WalletAddress: serverWalletAddress,
	}
}

type Transaction struct {
	con Contract
	to  common.Address
	Acc ServerAcc
}

func newTrans(to string, Con Contract, Acc ServerAcc) Transaction {
	recipientAddr := common.HexToAddress(to)
	return Transaction{
		con: Con,
		to:  recipientAddr,
		Acc: Acc,
	}
}

//get transaction details
func (tx *Transaction) getTxParams() {
	client := tx.con.client
	nonce, err := client.PendingNonceAt(context.Background(), tx.Acc.WalletAddress)
	if err != nil {
		log.Fatalf("error getting nonce %+v", err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("error on suggest gas price %+v", err)
	}
	auth := tx.Acc.auth
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
}

func (Tx *Transaction) MintNew() {
	CAToken := Tx.con.CAToken
	auth := Tx.Acc.auth
	to := Tx.to
	Tx.getTxParams()
	//sends the transcation
	//tx used for tracking the transcation detail
	transX, err := CAToken.Mint(auth, to)

	//add tx info on db list for tracking
}
