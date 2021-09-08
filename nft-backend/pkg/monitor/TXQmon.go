package monitor

import (
	"crypto/ecdsa"
	"errors"
	"sync"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

/*

user -> tx send

			ATOMIC
user -> [srvr -> db as pending q] -> [Qmon -> tx send -> block until mined or failure]

blockMon -> addsToDBForAccess

user w/ proof -> srvr -> validates against chain -> returns data

*/
type NewTxQ struct {
	workChan chan interface{}
}

type TXQmon struct {
	db  *redisDb.Database
	eth *blockchain.Ethereum
}

func NewQmon(redisAddr string, numWorkers int) *TXQmon {
	//new server instance
	tasks.NewServerClient(redisAddr, numWorkers)
	//New key manager instance

	Qmon := &TXQmon{}
	return Qmon
}

var ErrNoKeys error = errors.New("no privk available")
var ErrKeyConflict error = errors.New("privk added twice")

// PrivkManager releases keys to
type PrivkManager struct {
	sync.Mutex
	// stores all possibe private keys OR knows where to go get them
	// string key is the ether account
	consumedMap  map[string]bool
	availableMap map[string]bool
	masterSetMap map[string]ecdsa.PrivateKey
}

func (pm *PrivkManager) AddPrivk(privk ecdsa.PrivateKey) error {
	pm.Lock()
	defer pm.Unlock()
	addr := ethcrypto.PubkeyToAddress(privk.PublicKey)
	_, ok := pm.masterSetMap[addr]
	if ok {
		return ErrKeyConflict
	}
	pm.masterSetMap[addr] = privk
	pm.availableMap[addr] = true
	return nil
}

func (pm *PrivkManager) GetWithLock() (ecdsa.PrivateKey, func(), error) {
	pm.Lock()
	defer pm.Unlock()
	var privkAddr string
	for _privkAddr := range pm.availableMap {
		privkAddr = _privkAddr
		break
	}
	if privkAddr == "" {
		return ecdsa.PrivateKey{}, nil, ErrNoKeys
	}
	privk := pm.masterSetMap[privkAddr]
	delete(pm.availableMap, privkAddr)
	pm.consumedMap[privkAddr] = true
	return privk, pm.free(privkAddr), nil
}

func (pm *PrivkManager) free(privkAddr string) func() {
	return func() {
		//does it need mutex lock? each worker will have a unique key
		pm.Lock()
		defer pm.Unlock()
		delete(pm.consumedMap, privkAddr)
		pm.availableMap[privkAddr] = true
	}
}

//worker with access to private key manager
type TxWorker struct {
	pm *PrivkManager
	q  *TXQmon
}

func (txw *TxWorker) Run() error {
	//get the private key
	privk, free, err := txw.pm.GetWithLock()
	if err != nil {
		return err
	}
	defer free()
	//send the transaction
	go txw.loop(privk)
	return nil
}

//loop until
func (txw *TxWorker) loop(privk ecdsa.PrivateKey) {
	//send transaction
	for {

		//check if transaction went through
	}
}
