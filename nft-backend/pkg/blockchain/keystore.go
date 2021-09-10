package blockchain

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

func createNewAccKs() string {
	var password string
	keydir := "./tmp"
	ks := keystore.NewKeyStore(keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Printf("Enter a password: ")
	fmt.Scanf("%s", &password)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	//find the last file in the director
	//change the name
	accAddr := account.Address.Hex()
	originalpath := keydir + "/" + findLatestFile(keydir)
	newpath := keydir + "/" + accAddr
	os.Rename(originalpath, newpath)
	return accAddr
}

func StoreKs(password, privateKey string) error {
	keydir := "./tmp"
	ks := keystore.NewKeyStore(keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	/*
		fmt.Printf("Enter a password: ")
		fmt.Scanf("%s", &password)
		fmt.Printf("Enter your privateKey: ")
		fmt.Scanf("%s", &privateKey)*/

	//takes privatekey , returns *ecdsa.PrivateKey
	ECDSAprivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("error converting hex to ECDSA in StoreKs, %+v", err)
	}
	//convert PrivateKey string to a *ecdsa.PrivateKey
	account, err := ks.ImportECDSA(ECDSAprivateKey, password)
	if err != nil {
		log.Fatal(err)
	}
	accAddr := account.Address.Hex()
	originalpath := keydir + "/" + findLatestFile(keydir)
	newpath := keydir + "/" + accAddr
	os.Rename(originalpath, newpath)
	return nil
}

func importKs() accounts.Account {
	var accAddr string
	fmt.Printf("enter account address: ")
	fmt.Scanf("%s", &accAddr)
	file := "./tmp/" + accAddr
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
	return account
}
