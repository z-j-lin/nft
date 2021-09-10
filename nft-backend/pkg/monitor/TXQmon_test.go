package monitor

import (
	"log"
	"math/big"
	"testing"

	tasks "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

func TestNewQmon(t *testing.T) {
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
	NewQmon(redisAddr, 0, eth)
}
