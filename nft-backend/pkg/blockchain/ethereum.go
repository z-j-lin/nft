package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethereum struct {
	client   *ethclient.Client
	contract *Contract
	accounts map[common.Address]accounts.Account
	keystore *keystore.KeyStore
}

/*initializes a client to rpc. */
func NewEtherClient(rpcurl, contractAddress string) (*ethereum, error) {

	ethClient, err := ethclient.Dial(rpcurl)
	if err != nil {
		fmt.Errorf("errror @ NewEtherClient unable to dial RPC endpoint: %v", err)
		return nil, err
	}

	eth := &ethereum{
		client: ethClient,
	}

	eth.loadaccount()
	eth.loadpasscode()
	eth.contract = &Contract{eth: eth}

	return eth, nil

}
func (eth *ethereum) loadpasscode() {
	var passcode, address string
	passcodes := make(map[common.Address]string)
	fmt.Printf("enter address: ")
	fmt.Scanf("%s", &address)
	fmt.Printf("enter passcode: ")
	fmt.Scanf("%s", &passcode)
	passcodes[common.HexToAddress(address)] = passcode
}

func (eth *ethereum) loadaccount() {
	var dirPath string
	fmt.Printf("enter realtive key dir path: ")
	fmt.Scanf("%s", &dirPath)
	accounts := make(map[common.Address]accounts.Account)
	//for testing only: implement keystore for more secure account storage
	ks := keystore.NewKeyStore(dirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	for _, wallet := range ks.Wallets() {
		for _, account := range wallet.Accounts() {
			accounts[account.Address] = account
		}
	}
	eth.keystore = ks
	eth.accounts = accounts
}

type Transaction struct {
	con Contract
	to  common.Address
	Acc ServerAcc
}

func NewTransaction(to string, Con Contract, Acc ServerAcc) Transaction {
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

/*
func (Tx *Transaction) MintNew() {
	CAToken := Tx.con.CAToken
	auth := Tx.Acc.auth
	to := Tx.to
	Tx.getTxParams()
	//sends the transcation
	//tx used for tracking the transcation detail
	//transX, err := CAToken.Mint(auth, to)

	//add tx info on db list for tracking
}
*/
