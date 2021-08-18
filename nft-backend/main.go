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
type loginReq struct {
	//indicator to see if the account is logged in
	signature string
	data      string
	account   string
}
type loginRes struct {
	isloggedin bool
}

func verify(account string, data string, signature string) bool {
	//converting the pubkey from hex string to byte
	publicKeyBytes := crypto.FromECDSAPub(crypto.HexToECDSA(account))
	//taking signed message and converting it from string to byte
	signedMessage := []byte(signature)
	//convert data into byte array
	databyte := []byte(data)
	//hash original data
	hash := crypto.keccak256Hash(databyte)
	//extract the public key from the message
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signedMessage)
	if err != nil {
		log.Fatalf("unable to verify wallet")
		return false
	}
	fmt.Println("the recovered key:", sigPublicKey)
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
	session.Save(req, w)
}

//assigns a session ID
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methpds", "POST")
	//container for the login json data
	var loginreq loginReq
	err := json.NewDecoder(r.Body).Decode(&loginreq)
	if err != nil {
		log.Fatalf("unable to decode")
	}
	//gets a cookie
	session, _ := store.Get(r, "cookie-name")
	//authenticate
	if verify(loginreq.account, loginreq.data, loginreq.signature) {
		//if verification is true let user in
		session.Values["authenticated"] = true
		session.Save(r, w)
		//send back login acknoledgment
		var loginres loginRes
		loginres.isloggedin = true
		json.NewEncoder(w).Encode(&loginres)
	}

}

func main() {
	const rpcurl = "HTTP://127.0.0.1:9545"
	const contractAddress = "0x097063E71919E1C4af55F6468DF5295C76993bFb"
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/buy", login)
	http.HandleFunc("/load", login)
	http.ListenAndServe(":8080", nil)
}
