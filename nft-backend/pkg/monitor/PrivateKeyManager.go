package monitor

import (
	"crypto/ecdsa"
	"sync"
)

// PrivkManager releases keys to
type PrivkManager struct {
	sync.Mutex
	// stores all possibe private keys OR knows where to go get them
	// string key is the ether account
	consumedMap  map[string]bool
	availableMap map[string]bool
	masterSetMap map[string]ecdsa.PrivateKey
}

func NewPrivKManager() (*PrivkManager, error) {
	consumedMap := make(map[string]bool)
	availableMap := make(map[string]bool)
	masterSetMap := make(map[string]ecdsa.PrivateKey)
	//New key manager instance
	pm := &PrivkManager{
		consumedMap:  consumedMap,
		availableMap: availableMap,
		masterSetMap: masterSetMap,
	}
	return pm, nil
}

func (pm *PrivkManager) AddPrivk(addr string, privk ecdsa.PrivateKey) error {
	pm.Lock()
	defer pm.Unlock()
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
