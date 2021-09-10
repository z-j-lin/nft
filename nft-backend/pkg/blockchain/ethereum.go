package blockchain

import (
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	Client    *ethclient.Client
	chainID   *big.Int
	Contract  *Contract
	Accounts  []accounts.Account
	Keystore  *keystore.KeyStore
	Keys      map[common.Address]*keystore.Key
	Passcodes map[common.Address]string
}

/*initializes a client to rpc. */
func NewEtherClient(rpcurl, contractAddress string, chainID *big.Int) (*Ethereum, error) {
	ethClient, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Printf("error @ NewEtherClient unable to dial RPC endpoint: %v", err)
		return nil, err
	}
	keys := make(map[common.Address]*keystore.Key)
	eth := &Ethereum{
		Client:  ethClient,
		chainID: chainID,
		Keys:    keys,
	}
	eth.loadaccount()
	eth.loadpasscode()
	err = eth.unlockkey()
	if err != nil {
		return nil, err
	}
	eth.Contract = &Contract{
		eth:             eth,
		ContractAddress: common.HexToAddress(contractAddress),
	}
	eth.Contract.init()
	return eth, nil
}
func (eth *Ethereum) loadpasscode() {
	passcodes := make(map[common.Address]string)
	passcode := "pineapple"
	for _, account := range eth.Accounts {
		address := account.Address
		passcodes[address] = passcode
	}
	eth.Passcodes = passcodes
}

func (eth *Ethereum) loadaccount() {
	//fmt.Printf("enter realtive key dir path: ")
	//fmt.Scanf("%s", &dirPath)
	dirPath := "./tmp"
	//for testing only: implement keystore for more secure account storage
	ks := keystore.NewKeyStore(dirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	eth.Keystore = ks
	eth.Accounts = ks.Accounts()
}

func (eth *Ethereum) unlockkey() error {
	for _, account := range eth.Accounts {
		passcode, exists := eth.Passcodes[account.Address]
		if exists {
			//get the encrypted private key in json form
			encrytpedKey, err := ioutil.ReadFile((account.URL.Path))
			if err != nil {
				log.Println("unable to getch key from URL path for", account.Address, "err:", err)
			}
			//decrypt the private key
			privateKey, err := keystore.DecryptKey(encrytpedKey, passcode)
			if err != nil {
				log.Fatal(err)
			}
			eth.Keys[account.Address] = privateKey
		}
	}
	return nil
}
