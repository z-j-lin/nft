package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

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
	SignedMessage string `json: "signedmessage"`
	AccountAddr   string `json: "accountaddr`
}
type loginRes struct {
	isloggedin bool
}

//struct to hold information about each token
type TokenRegistry struct {
	TokenID   uint
	Account   string
	StartTime time.Time
	EndTime   time.Time
}

//slice to hold information about all tokens

//create a map that mapps sessionIDs to etheraddress
//create a queue for sending transactions
//need to monitor if transactions go through

//handler function for buying a token
func BuyToken(w http.ResponseWriter, r *http.Request) {
	//verify user credential
	//if user is verfied pull out transaction information
	//run buildtransaction with inforrmation
	//add transaction to the queue
}

func BuildTransactionMint() {
}
func SendTransactionMint() {

}

func VerifyTransactionMint() {
}

//sends the client a list of tokens owned by the address and associated resource ID
func LoadAccessTokens(w http.ResponseWriter, r *http.Request) {
	//decode request body
	//extract session ID, user account

	//encode new body, sends back array tokens
}

func fetchResource(w http.ResponseWriter, r *http.Request) {

}

//cleanup function needs to be ran as a go routine

//verfication for metamask login message
func verify(account string, data string, signature string) bool {
	/*
			   PUT DATA HERE
					{
		    "signedmessage": "0x58c268c7e3fdf11e13fe9e05f612e4d44b28a333a55630c551e04aef633f6d2825b790163f98632a66a6745d9ed8f0430785f3a7d997e28b595ce388018ce01f1c",
		    "accountaddr" : "0x28cB37028ECE65435480565c7f71f8a372bb655d"
		}
	*/
	fmt.Println(data)
	//converting the pubkey from hex string to byte
	//taking signed message and converting it from string to byte
	signedMessage, err := hex.DecodeString(signature[2:])
	if err != nil {
		panic(err)
	}
	signedMessage[64] -= 27
	AccountAddr, err := hex.DecodeString(account[2:])
	fmt.Println("string decoded AccAddr[2:]: ", hex.EncodeToString(AccountAddr))
	if err != nil {
		panic(err)
	}
	validationMsg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(data)) + data
	//convert data into byte array

	databyte := []byte(validationMsg)
	//hash original data
	hash := crypto.Keccak256Hash(databyte)
	//extract the public key from the message

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signedMessage)
	if err != nil {
		panic(err)
	}
	fmt.Println("sigPubkey:", sigPublicKey)
	fmt.Println("accoutnaddr:", AccountAddr)
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		panic(err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	fmt.Println("the recovered address:", recoveredAddr)
	fmt.Println("the address:", account)
	if err != nil {
		log.Fatalf("unable to verify wallet: %v %v, signature: %x", err, len(signature), signedMessage)
		return false
	}
	//check if the recovered address matches the actual address
	/*matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	if matches {
		return true
	} else {
		return false
	}*/
	return false
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

//handler for login
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	//container for the login json data
	var lr loginReq
	fmt.Println(1)
	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		log.Fatalf("unable to decode biatch %v", err)
	}

	//gets a cookie
	session, _ := store.Get(r, "cookie-name")
	//authenticate
	if verify(lr.AccountAddr, "hello", lr.SignedMessage) {
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
	signedmessage := "0x58c268c7e3fdf11e13fe9e05f612e4d44b28a333a55630c551e04aef633f6d2825b790163f98632a66a6745d9ed8f0430785f3a7d997e28b595ce388018ce01f1c"
	accountaddr := "0x28cB37028ECE65435480565c7f71f8a372bb655d"
	if ok := verify(accountaddr, "hello", signedmessage); !ok {
		panic("not verified")
	}
	return
	//TODO: ZJ REMOVE ABOVE THIS LINE
	var dir string
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	const rpcurl = "HTTP://127.0.0.1:9545"
	//contractAddress := "0x097063E71919E1C4af55F6468DF5295C76993bFb"
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/buy", BuyToken).Methods("POST")
	//sends back an array resources owned by address
	router.HandleFunc("/load", LoadAccessTokens).Methods("GET")
	router.HandleFunc("request/{resourceID}", fetchResource).Methods("GET")
	http.Handle("/", router)

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
