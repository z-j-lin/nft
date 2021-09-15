package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//connection variables
const (
	rpcurl    = "HTTP://127.0.0.1:9545"
	redisAddr = "127.0.0.1:6379"
)

/*
	isloggedin bool
	signature string
	data string
	publicKeyECDSA string
*/

type loginRes struct {
	Isloggedin string
}

//struct to hold information about each token
type TokenRegistry struct {
	resourceID string
	Account    string
}

func ContentStore(w http.ResponseWriter, r *http.Request) {
	log.Println("reached the login endpoint")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	db, err := redisDb.NewDBinstance()
	if err != nil {
		panic(err)
	}
	Cstore, err := db.GetStore()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	//encode the object into byte form
	respBytes, err := json.Marshal(Cstore)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	} else {
		//send back the store content
		w.Write(respBytes)
	}

}

//handler function for buying a token
func BuyToken(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		log.Println("reached the buytoken endpoint preflight")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.WriteHeader(http.StatusOK)
	}
	if r.Method == "POST" {
		log.Println("reached the buytoken endpoint")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		//get the content of type of the incoming request
		contentT := r.Header.Get("Content-Type")
		if contentT == "application/json" {
			//container for the login json data
			var buy TokenRegistry
			//decode body
			err := json.NewDecoder(r.Body).Decode(&buy)
			if err != nil {
				log.Println("unable to decode json pack", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened!"))
				return
			}
			log.Println("buy account:", buy.Account)
			session, _ := store.Get(r, buy.Account)
			//verify user credential
			auth, ok := session.Values["authenticated"].(bool)
			log.Println("authentication", auth, ok)
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			TC := tasks.NewTaskClient(redisAddr)

			//add transaction to the queue
			// TODO: HANDLE ERROR AND PASS INFO BACK TO CALLER
			err = TC.QMintTask(buy.Account, buy.resourceID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened!"))
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
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
	//check against block chain if client owns token

}

//verfication for metamask login message
func verify(account string, data string, signature string) (bool, error) {
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

type loginReq struct {
	//indicator to see if the account is logged in
	Signature string `json: "signature"`
	Account   string `json: "account"`
}

//handler for login
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		log.Println("reached the login endpoint preflight")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.WriteHeader(http.StatusOK)
	}
	if r.Method == "POST" {
		log.Println("reached the login endpoint")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		//get the content of type of the incoming request
		contentT := r.Header.Get("Content-Type")
		if contentT == "application/json" {
			//container for the login json data
			var loginR loginReq
			err := json.NewDecoder(r.Body).Decode(&loginR)
			if err != nil {
				log.Println("unable to decode json pack", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened!"))
				return
			}
			//gets a cookie
			log.Println("account logging in:", loginR.Account)
			session, _ := store.Get(r, loginR.Account)
			//authenticate
			isOwner, err := verify(loginR.Account, "hello", loginR.Signature)
			if err != nil {
				panic(err)
			}
			//if verification is true let user in
			if isOwner {
				//TODO: set session ID
				log.Println("logged in:", session.Values["authenticated"])
				session.Values["authenticated"] = true
				log.Println("authenticated in:", session.Values["authenticated"])
				auth, ok := session.Values["authenticated"].(bool)
				log.Println("authentication", auth, ok)
				session.Save(r, w)
				auth, ok = session.Values["authenticated"].(bool)
				log.Println("authentication", auth, ok)
				log.Printf("session ID: %s, isNew: %t, name: %s", session.ID, session.IsNew, session.Name())
				//send back login acknoledgment
				loginres := loginRes{
					Isloggedin: "1",
				}
				respbyte, err := json.Marshal(loginres)
				if err != nil {
					panic(err)
				}
				w.Write(respbyte)
			}
		}
	}
}

func main() {
	var dir string
	store.Options.Path = "/"
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")

	//contractAddress := "0x097063E71919E1C4af55F6468DF5295C76993bFb"
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/buy", BuyToken).Methods("POST", "OPTIONS")
	//sends back an array resources owned by address
	router.HandleFunc("/load", LoadAccessTokens).Methods("GET")
	router.HandleFunc("/getstore", ContentStore).Methods("GET")
	router.HandleFunc("request/{resourceID}", fetchResource).Methods("GET")
	http.Handle("/", router)

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
