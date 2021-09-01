package monitor

import (
	"testing"

	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

func TesteventlogByTXHash(t *testing.T) {

	db, err := redisDb.NewDBinstance()
	if err != nil {
		t.Error(err)
	}
	//send out a transaction

	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	contractAddr := "0x3DA85558aF6D0d0D03283fa23eD1edE90f7E3E03"
	eth, err := blockchain.NewEtherClient(rpcurl, contractAddr)
	Events := Events{
		eth: eth,
	}

	Events.eventlogByTXHash()
}
