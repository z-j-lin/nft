package monitor

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func subToConEventLog() {
	rpcurl := "wss://ropsten.infura.io/ws/v3/27c2937f16d14d33a4c8315e22109f09"
	client, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Panic("cant connect to rpc server: %+v", err)
	}
	contractAddress := common.HexToAddress("0x3DA85558aF6D0d0D03283fa23eD1edE90f7E3E03")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Panic("unable to filter logs:", err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog)
		}
	}
}
