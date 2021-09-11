package monitor

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

func TestNewQmon(t *testing.T) {
	rdb, err := redisDb.NewDBinstance()
	if err != nil {
		panic(err)
	}
	chainID := big.NewInt(int64(5777))
	eth, err := blockchain.NewEtherClient("HTTP://127.0.0.1:9545", "0x810dA0c61C3b19087d40cdCa990790351F146dc8", chainID)
	if err != nil {
		log.Fatal(err)
	}
	//load tasks
	redisAddr := "127.0.0.1:6379"
	TC := tasks.NewAsyncredisClient(redisAddr)
	var tests = []struct {
		address   string
		contentID string
	}{
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing2"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing3"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing4"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing5"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing6"},
	}
	for _, test := range tests {
		err := TC.QMintTask(test.address, test.contentID)
		if err != nil {
			t.Error(err)
		}
	}
	//start server to process tasks
	NewQmon(redisAddr, 0, eth, rdb)
}

func TestAddPrivk(t *testing.T) {
	chainID := big.NewInt(int64(5777))
	eth, err := blockchain.NewEtherClient("HTTP://127.0.0.1:9545", "0x810dA0c61C3b19087d40cdCa990790351F146dc8", chainID)
	if err != nil {
		log.Fatal(err)
	}
	consumedMap := make(map[string]bool)
	availableMap := make(map[string]bool)
	masterSetMap := make(map[string]ecdsa.PrivateKey)
	//New key manager instance
	pm := &PrivkManager{
		consumedMap:  consumedMap,
		availableMap: availableMap,
		masterSetMap: masterSetMap,
	}
	//add keys to the key manager
	for addr, key := range eth.Keys {
		key := *key.PrivateKey
		pm.AddPrivk(addr.Hex(), key)
		log.Println(addr, key)
	}

}

func TestRunTXWorker(t *testing.T) {
	db, err := redisDb.NewDBinstance()
	if err != nil {
		log.Fatal(err)
	}
	chainID := big.NewInt(int64(5777))
	eth, err := blockchain.NewEtherClient("HTTP://127.0.0.1:9545", "0x810dA0c61C3b19087d40cdCa990790351F146dc8", chainID)
	if err != nil {
		log.Fatal(err)
	}
	consumedMap := make(map[string]bool)
	availableMap := make(map[string]bool)
	masterSetMap := make(map[string]ecdsa.PrivateKey)
	//New key manager instance
	pm := &PrivkManager{
		consumedMap:  consumedMap,
		availableMap: availableMap,
		masterSetMap: masterSetMap,
	}
	//add keys to the key manager
	for addr, key := range eth.Keys {
		pm.AddPrivk(addr.Hex(), *key.PrivateKey)
	}
	//used in the server client handler functions
	hdl := &Handler{
		PrivkManager: pm,
		eth:          eth,
		db:           db,
	}
	AccountAddress, ResourceID := "0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"
	send := blockchain.NewMintTransaction(AccountAddress, ResourceID, hdl.eth, hdl.db)
	//start a new txworker, with the mint transaction function
	NewTX := TxWorker{
		pm:     hdl.PrivkManager,
		eth:    hdl.eth,
		sendTX: send,
	}
	err = NewTX.Run()
	//returns the status of the job
	if err != nil {
		t.Error(err)
	}
}
