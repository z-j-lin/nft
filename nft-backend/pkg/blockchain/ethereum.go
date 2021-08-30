package blockchain

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	client    *ethclient.Client
	contract  *Contract
	account   accounts.Account
	keystore  *keystore.KeyStore
	key       *keystore.Key
	passcodes map[common.Address]string
}

/*initializes a client to rpc. */
func NewEtherClient(rpcurl, contractAddress string) (*Ethereum, error) {

	ethClient, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Printf("errror @ NewEtherClient unable to dial RPC endpoint: %v", err)
		return nil, err
	}

	eth := &Ethereum{
		client: ethClient,
	}

	eth.loadaccount()
	eth.loadpasscode()
	err = eth.unlockkey(eth.account)
	if err != nil {
		return nil, err
	}
	eth.contract = &Contract{eth: eth}
	eth.contract.init()

	return eth, nil

}
func (eth *Ethereum) loadpasscode() {
	var passcode, address string
	passcodes := make(map[common.Address]string)
	fmt.Printf("enter address: ")
	fmt.Scanf("%s", &address)
	fmt.Printf("enter passcode: ")
	fmt.Scanf("%s", &passcode)
	passcodes[common.HexToAddress(address)] = passcode
	eth.passcodes = passcodes
}

func (eth *Ethereum) loadaccount() {
	var dirPath string
	fmt.Printf("enter realtive key dir path: ")
	fmt.Scanf("%s", &dirPath)
	//for testing only: implement keystore for more secure account storage
	ks := keystore.NewKeyStore(dirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	eth.keystore = ks
	eth.account = ks.Accounts()[0]
}

func (eth *Ethereum) unlockkey(account accounts.Account) error {
	passcode, exists := eth.passcodes[account.Address]
	if !exists {
		return fmt.Errorf("passcode not found")
	}
	//get the encrypted private key in json form
	encrytpedKey, err := ioutil.ReadFile((account.URL.Path))
	if err != nil {
		return err
	}
	//decrypt the private key
	privateKey, err := keystore.DecryptKey(encrytpedKey, passcode)
	eth.key = privateKey
	return nil
}
