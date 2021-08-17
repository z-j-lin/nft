package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/sessions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//verfication for metamask login message
/*
func verification(publicKeyECDSA string, data, signature byte) bool {
	publicKeyBytes := crypto.FromECDSAPub(crypto.HexToECSDA(publicKeyECDSA))
	hash := crypto.keccak256Hash(data)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err {
		log.Fatalf("unable to verify wallet")
		return false
	}
	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	return true
}
*/
var (
	key   = []byte("mulva")
	store = sessions.NewCookieStore(key)
)

type loginR struct {
	isLoggedIn bool
}

func EtherInit(rpcurl, contractAddress string) {
	client, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Fatalf("Failed to connect to the ether network: %v", err)
	}
	fmt.Println("we are connected")
	//contract address
	conAddress := common.HexToAddress(contractAddress)
	fmt.Println("contract address: " + conAddress.Hex())
	accountBal, err := client.BalanceAt(context.Background(), conAddress, nil)
	fmt.Println("Account Balance:", accountBal)
	//private key needs a keystore
	privateKey, err := crypto.HexToECDSA("9846163dfc41a7d467f6c35c40e24408d972db8d30f3c886adadbcb341f58c6e")
	if err != nil {
		log.Fatalf("private key problem: %v", err)
	}
	publicKey := privateKey.Public()
	fmt.Println(publicKey)
}
func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	rData := loginR{
		isLoggedIn: true,
	}
	session.Values["authenticated"] = true
	w.Header().Set("Content-Type", "application/json")

	Data, err := json.Marshal(rData)
	if err != nil {
		log.Fatalf("JSON problem server side: %v", err)
	}
	w.Write(Data)
	session.Save(r, w)
}

func main() {
	const rpcurl = "HTTP://127.0.0.1:9545"
	const contractAddress = "0x097063E71919E1C4af55F6468DF5295C76993bFb"
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}
