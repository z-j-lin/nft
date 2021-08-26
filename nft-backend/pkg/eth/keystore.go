package keymanager

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

type Keys struct {
	keydir  string //./tmp
	account accounts.Account
}

func findLatestFile(dir string) string {
	files, err := ioutil.ReadDir(dir + "/.")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	var modTime time.Time
	var mostRecentFi string
	for _, fi := range files {
		fmt.Println(fi.Name(), fi.ModTime())
		if newModtime := fi.ModTime(); newModtime.After(modTime) {
			modTime = fi.ModTime()
			mostRecentFi = fi.Name()
		}
	}
	return mostRecentFi
}

func (k *Keys) createNewAccKs() string {
	var password string
	ks := keystore.NewKeyStore(k.keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Printf("Enter a password: ")
	fmt.Scanf("%s", &password)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	//find the last file in the director
	//change the name
	accAddr := account.Address.Hex()
	originalpath := k.keydir + "/" + findLatestFile(k.keydir)
	newpath := k.keydir + "/" + accAddr
	os.Rename(originalpath, newpath)
	return accAddr
}

func (k *Keys) StoreKs() (string, error) {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Println()
	var (
		password, privateKey string
	)
	fmt.Printf("Enter a password: ")
	fmt.Scanf("%s", &password)
	fmt.Printf("Enter your privateKey: ")
	fmt.Scanf("%s", &privateKey)
	//takes privatekey string removes 0x, returns *ecdsa.PrivateKey
	ECDSAprivateKey, err := crypto.HexToECDSA(privateKey[2:])
	if err != nil {
		log.Fatalf("error converting hex to ECDSA in StoreKs, %+v", err)
	}
	//convert PrivateKey string to a *ecdsa.PrivateKey
	account, err := ks.ImportECDSA(ECDSAprivateKey, password)
	k.account = account
	accAddr := account.Address.Hex()
	originalpath := k.keydir + "/" + findLatestFile(k.keydir)
	newpath := k.keydir + "/" + accAddr
	os.Rename(originalpath, newpath)
	return accAddr, nil
}

func importKs() {
	file := "./tmp/.*"
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var password string
	fmt.Scanf("%s", &password)
	account, err := ks.Import(jsonBytes, password, password)
	fmt.Println(account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//createKs()
	importKs()
}
