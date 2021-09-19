package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

var (
	key      = []byte("super-secret-key")
	store    = sessions.NewCookieStore(key)
	chainID  = big.NewInt(int64(5444))
	eth      *blockchain.Ethereum
	db       *redisDb.Database
	TC       *tasks.TaskClient
	NonceMan *tasks.NonceMan
)

//connection variables(should be in a config file)
const (
	rpcurl    = "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	redisAddr = "127.0.0.1:6379"

	contractAddress = "0xb410756d52b1250aB9bE358437Ab41a4D7636Af8"
)

type loginRes struct {
	Isloggedin string
}

//struct to hold information about each token
type TokenRegistry struct {
	Account    string `json:"account"`
	ResourceID string `json:"resourceid"`
}

func ContentStore(w http.ResponseWriter, r *http.Request) {
	log.Println("reached the login endpoint")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

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
func PostHeaders(w http.ResponseWriter, r *http.Request) {
	log.Println("reached the buytoken endpoint preflight")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

//handler function for buying a token
func BuyToken(w http.ResponseWriter, r *http.Request) {
	PostHeaders(w, r)
	log.Println("reached the buytoken endpoint")
	//get the content of type of the incoming request
	contentT := r.Header.Get("Content-Type")
	if r.Method == "POST" && contentT == "application/json" {
		session, _ := store.Get(r, "access")
		//check if the user is logged in
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
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

		//add transaction to the queue
		nonce := NonceMan.GetnonceWithLock()
		err = TC.QMintTask(buy.Account, buy.ResourceID, nonce)
		log.Printf("RESTAPI: Mint Transaction task queued. Account: %s, ResourceID: %s, Nonce: %d", buy.Account, buy.ResourceID, nonce)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

type contentReq struct {
	Account string `json:"account"`
	TokenID string `json:"tokenid"`
}

//on request send back access tokens to user
func LoadAccessTokens(w http.ResponseWriter, r *http.Request) {
	PostHeaders(w, r)
	log.Println("reached the access tokens endpoint")
	contentT := r.Header.Get("Content-Type")
	if r.Method == "POST" && contentT == "application/json" {
		session, _ := store.Get(r, "access")
		//check if the user is logged in
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var conR contentReq
		//extract the user account
		err := json.NewDecoder(r.Body).Decode(&conR)
		if err != nil {
			log.Println("unable to decode json pack", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}
		//query redis db for tokens owned by account
		AccTokens, err := db.GetAccTokens(conR.Account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}
		//encode the object into byte form
		respBytes, err := json.Marshal(AccTokens)
		if err != nil {
			log.Println("loadacccesstokens:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
		} else {
			//send back array of tokens owned by client
			w.WriteHeader(http.StatusOK)
			w.Write(respBytes)
		}
	}
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

func logout(w http.ResponseWriter, r *http.Request) {
	PostHeaders(w, r)
	session, _ := store.Get(r, "access")
	//revoke permission
	session.Values["authenticated"] = false
	session.Save(r, w)
	loginres := loginRes{
		Isloggedin: "0",
	}
	respbyte, err := json.Marshal(loginres)
	if err != nil {
		panic(err)
	}
	w.Write(respbyte)
}

type loginReq struct {
	//indicator to see if the account is logged in
	Signature string `json:"signature"`
	Account   string `json:"account"`
}

//handler for login
func login(w http.ResponseWriter, r *http.Request) {
	PostHeaders(w, r)
	if r.Method == "POST" {
		log.Println("reached the login endpoint")
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
			session, _ := store.Get(r, "access")
			session.Options.Path = "/"
			session.Options.HttpOnly = true
			session.Options.SameSite = http.SameSiteNoneMode
			session.Options.Secure = true
			//authenticate
			isOwner, err := verify(loginR.Account, "hello", loginR.Signature)
			if err != nil {
				panic(err)
			}
			//if verification is true let user in
			if isOwner {
				session.Values["authenticated"] = true
				log.Println("authenticated in:", session.Values["authenticated"])
				session.Save(r, w)
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
func fetchResource(w http.ResponseWriter, r *http.Request) {
	PostHeaders(w, r)
	log.Println("reached the access tokens endpoint")
	contentT := r.Header.Get("Content-Type")
	if r.Method == "POST" && contentT == "application/json" {
		session, _ := store.Get(r, "access")
		//check if the user is logged in
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		//container for request parameters: account addr, resourceID
		var ResourceReq contentReq
		//extract the user account
		err := json.NewDecoder(r.Body).Decode(&ResourceReq)
		if err != nil {
			log.Println("unable to decode json pack", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}
		//check if the requester is the owner of the token
		isowner, err := eth.Contract.IsOwner(ResourceReq.Account, ResourceReq.TokenID)
		if isowner && err == nil {
			ResourceID, err := eth.Contract.GetResourceID(ResourceReq.TokenID)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			filebyte, err := ioutil.ReadFile(ResourceID + ".jpeg")
			if err != nil {
				log.Println("Fetresource:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened!"))
			} else {
				//send back array of tokens owned by client
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "aplication/json")
				w.Write([]byte(ResourceID))
				w.Header().Set("Content-Type", "aplication/octet-stream")
				w.Write(filebyte)
			}
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
}
func main() {

	var dir string
	ethC, err := blockchain.NewKeylessEthClient(rpcurl, contractAddress, chainID)
	if err != nil {
		panic(err)
	}
	eth = ethC
	eth.Contract = blockchain.NewContract(eth, contractAddress)
	TC = tasks.NewTaskClient(eth, redisAddr)
	NonceMan = tasks.NewNonceManager(eth)
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	rdb, err := redisDb.NewDBinstance()
	if err != nil {
		panic(err)
	}
	db = rdb
	//contractAddress := "0x097063E71919E1C4af55F6468DF5295C76993bFb"
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", logout).Methods("POST", "OPTIONS")

	router.HandleFunc("/buy", BuyToken).Methods("POST", "OPTIONS")
	//sends back an array of resources owned by address
	router.HandleFunc("/load", LoadAccessTokens).Methods("POST", "OPTIONS")
	router.HandleFunc("/getstore", ContentStore).Methods("GET")

	router.HandleFunc("/request", fetchResource).Methods("POST", "OPTIONS")
	http.Handle("/", router)

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
