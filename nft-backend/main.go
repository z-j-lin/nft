package main

import (
	"bytes"
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
var (
	key   = []byte("mulva")
	store = sessions.NewCookieStore(key)
)

/*
	isloggedin bool
	signature string
	data string
	publicKeyECDSA string
*/
type loginR struct {
	//indicator to see if the account is logged in
	isloggedin     bool
	signature      string
	data           string
	publicKeyECDSA string
}

func verify(publicKeyECDSA string, data string, signature string) bool {
	//converting the pubkey from hex string to byte
	publicKeyBytes := crypto.FromECDSAPub(crypto.HexToECSDA(publicKeyECDSA))
	//taking signed message and converting it from string to byte
	signedMessage := []byte(signature)
	//convert data into byte array
	databyte := []byte(data)
	//hash original data
	hash := crypto.keccak256Hash(databyte)
	//extract the public key from the message
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signedMessage)
	if err {
		log.Fatalf("unable to verify wallet")
		return false
	}
	//check if it matches
	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	if matches {
		return true
	} else {
		return false
	}
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

func logout(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	//revoke permission
	session.Values["authenticated"] = false
	session.save(req, w)
}

//assigns a session ID
func login(w http.ResponseWriter, r *http.Request) {
	//container for the login json data
	var loginreq loginR
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginreq)
	if err != nil {
		log.Fatalf("unable to decode")
	}
	//gets a cookie
	session, _ := store.Get(r, "cookie-name")
	//authenticate
	if verify(loginreq.publicKeyECDSA, loginreq.data, loginreq.signature) {
		//if verification is true let user in
		session.Values["authenticated"] = true
		w.Header().Set("Content-Type", "application/json")

		Data, err := json.Marshal(rData)
		if err != nil {
			log.Fatalf("JSON problem server side: %v", err)
		}
		session.Save(r, w)
	}
	w.Write(Data)

}

func main() {
	const rpcurl = "HTTP://127.0.0.1:9545"
	const contractAddress = "0x097063E71919E1C4af55F6468DF5295C76993bFb"
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}
