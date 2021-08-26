package main

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
	rpcurl := "https://ropsten.infura.io/v3/27c2937f16d14d33a4c8315e22109f09"
	client, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Panic("cant connect to rpc server: %+v", err)
	}
	contractAddress := common.HexToAddress("5DD32251e8A2Ac4D40d115520dC49E6a94465ee7")
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

func main() {
	subToConEventLog()

}
