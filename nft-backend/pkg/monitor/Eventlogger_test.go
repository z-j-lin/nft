package monitor

import (
	"testing"

	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

func TestEventlogByTXHash(t *testing.T) {
	/*db, err := redisDb.NewDBinstance()
	if err != nil {
		t.Error(err)
	}*/
	//taskStatus := make(chan bool)
	//send out a transaction
	//toAddr := "0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f"
	//resourceID := "thing1"
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	contractAddr := "0x3DA85558aF6D0d0D03283fa23eD1edE90f7E3E03"
	eth, err := blockchain.NewEtherClient(rpcurl, contractAddr)
	if err != nil {
		t.Error(err)
	}
	//blockchain.NewTransaction(toAddr, resourceID, eth.Contract, db, taskStatus)
	//done := <-taskStatus
	//txHashHex, accountaddrHex, rID := db.DQpending()
	txHashHex := "0xfd7c865a20749be0fbd05aae462989d10e37ef288d6688be02dd64904ed43766"
	Events := Events{
		eth: eth,
	}

	//_ = accountaddrHex
	//_ = rID
	//_=done

	Events.EventlogByTXHash(txHashHex)
}
