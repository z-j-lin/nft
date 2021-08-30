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
	Client    *ethclient.Client
	Contract  *Contract
	Account   accounts.Account
	Keystore  *keystore.KeyStore
	Key       *keystore.Key
	Passcodes map[common.Address]string
}

/*initializes a client to rpc. */
func NewEtherClient(rpcurl, contractAddress string) (*Ethereum, error) {

	ethClient, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Printf("errror @ NewEtherClient unable to dial RPC endpoint: %v", err)
		return nil, err
	}

	eth := &Ethereum{
		Client: ethClient,
	}

	eth.loadaccount()
	eth.loadpasscode()
	err = eth.unlockkey(eth.Account)
	if err != nil {
		return nil, err
	}
	eth.Contract = &Contract{eth: eth}
	eth.Contract.init()
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
	eth.Passcodes = passcodes
}

func (eth *Ethereum) loadaccount() {
	var dirPath string
	fmt.Printf("enter realtive key dir path: ")
	fmt.Scanf("%s", &dirPath)
	//for testing only: implement keystore for more secure account storage
	ks := keystore.NewKeyStore(dirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	eth.Keystore = ks
	eth.Account = ks.Accounts()[0]
}

func (eth *Ethereum) unlockkey(account accounts.Account) error {
	passcode, exists := eth.Passcodes[account.Address]
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
	eth.Key = privateKey
	return nil
}
