package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	redisDB "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
)

var (
	key   = []byte("super-secret-key")
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
	isloggedin string
}

//struct to hold information about each token
type TokenRegistry struct {
	resourceID string
	Account    string
}

//create a map that mapps sessionIDs to etheraddress

//create a queue for sending transactions

//need to monitor if transactions go through

//handler function for buying a token
func BuyToken(w http.ResponseWriter, r *http.Request) {
	//instantiate db instance
	rdb, err := redisDB.NewDBinstance()
	if err != nil {
		panic(err)
	}
	var buy TokenRegistry
	//decode body
	err = json.NewDecoder(r.Body).Decode(&buy)
	if err != nil {
		panic(err)
	}
	session, _ := store.Get(r, buy.Account)
	//verify user credential
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	//add transaction to the queue
	rdb.Qmint(buy.Account, buy.resourceID)
}

//on request send back access tokens to user
func LoadAccessTokens(w http.ResponseWriter, r *http.Request) {
	//extract the user account
	//querrry redis db for tokens owned by account
	//query each tokenID for the

	//extract session ID, user account

	//encode new body, sends back array tokens
}

func fetchResource(w http.ResponseWriter, r *http.Request) {

}

//cleanup function needs to be ran as a go routine

//verfication for metamask login message
func verify(account string, data string, signature string) (bool, error) {

	fmt.Println(data)
	//takes signed message and convert it from string to byte array
	signedMessage, err := hex.DecodeString(signature[2:])
	if err != nil {
		panic(err)
	}
	//set this to indicate a ethereum signed message
	signedMessage[64] -= 27
	if err != nil {
		panic(err)
	}
	//concatenate the data header, string length, and data
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
	//returns the recovered ECDSA pubkey
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		panic(err)
	}
	//extracts the ECDSA pubkey as a hex string
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	accByte, err := hex.DecodeString(account[2:])
	if err != nil {
		panic(err)
	}
	//the byte array of the recovered address
	recAccByte := recoveredAddr.Bytes()
	//check if the recovered address matches the actual address
	matches := bytes.Equal(accByte, recAccByte)
	if matches {
		return true, nil
	} else {
		return false, nil
	}
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
	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		log.Fatalf("unable to decode biatch %v", err)
	}

	//gets a cookie
	session, _ := store.Get(r, lr.AccountAddr)
	//authenticate
	isOwner, err := verify(lr.AccountAddr, "hello", lr.SignedMessage)
	if err != nil {
		panic(err)
	}
	//if verification is true let user in
	if isOwner {
		//TODO: set session ID
		session.Values["authenticated"] = true
		session.Save(r, w)
		log.Printf("session ID: %s, isNew: %t, name: %s", session.ID, session.IsNew, session.Name())
		//send back login acknoledgment
		loginres := loginRes{
			isloggedin: "1",
		}
		//TODO: fix issue with respose body
		respbyte, err := json.Marshal(loginres)
		if err != nil {
			panic(err)
		}
		//TODO: DELETE
		fmt.Println(respbyte)
		w.Write([]byte("hello"))
	}
}

func main() {

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
