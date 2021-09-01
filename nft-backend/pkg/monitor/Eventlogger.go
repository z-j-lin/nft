package monitor

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type LogMinted struct {
	To      common.Address
	tokenID *big.Int
}

type LogDeletedTokens struct {
	deleteIds []*big.Int
}

type Events struct {
	eth *blockchain.Ethereum
}

func (e *Events) eventlogByTXHash(txhex string) {
	txhash := common.HexToHash(txhex)
	receipt, err := e.eth.Client.TransactionReceipt(context.TODO(), txhash)
	if err != nil {
		log.Println(err, "@ eventlogByTXHASH")
	}
	//gete the logs of the transaction
	TXlog := receipt.Logs

	CATokenAbi, err := abi.JSON(strings.NewReader(CAToken.CATokenABI))
	if err != nil {
		log.Println(err, "@ catokenABI")
	}

	LogMintedSig := []byte("Minted(address,uint256)")
	LogDeletedTokensSig := []byte("DeletedTokens(uint256[])")
	LogMintedSigHash := crypto.Keccak256Hash(LogMintedSig)
	LogDeletedTokensSigHash := crypto.Keccak256Hash(LogDeletedTokensSig)
	for _, vLog := range TXlog {

		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case LogMintedSigHash.Hex():
			fmt.Printf("Log Name: Minted\n")
			MintedEvent, err := CATokenAbi.Unpack("Minted", vLog.Data)
			if err != nil {
				log.Println(err, "@ catokenABI")
			}

			fmt.Println(MintedEvent...)
		case LogDeletedTokensSigHash.Hex():
			fmt.Printf("Log Name: DeletedTokens\n")
			DeletedTokenEvent, err := CATokenAbi.Unpack("DeletedTokens", vLog.Data)
			if err != nil {
				log.Println(err, "@ catokenABI")
			}
			fmt.Println(DeletedTokenEvent...)
		}
	}
}

/*
func (e *Events) eventloggerbyBlock(txhex string, from, to int64) {
	ConAddr := e.eth.Contract.ContractAddress
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(to),
		Addresses: []common.Address{
			ConAddr,
		},
	}
	logs, err := e.eth.Client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Print(err, "at eventlo0gger fileterlogs")
	}
	CATokenAbi, err := abi.JSON(strings.NewReader(CAToken.CATokenABI))
	if err != nil {
		log.Println(err, "@ catokenABI")
	}
	LogMintedSig := []byte("Minted(address,uint256)")
	LogDeletedTokensSig := []byte("DeletedTokens(uint256[])")
	LogMintedSigHash := crypto.Keccak256Hash(LogMintedSig)
	LogDeletedTokensSigHash := crypto.Keccak256Hash(LogDeletedTokensSig)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case LogMintedSigHash.Hex():
			fmt.Printf("Log Name: Minted\n")

			var mintEvent LogMinted

			 := CATokenAbi.Unpack("Minted", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			mintEvent.to = common.HexToAddress(vLog.Topics[1].Hex())
			mintEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())

		case LogDeletedTokensSigHash.Hex():
			fmt.Printf("Log Name: Approval\n")


			err := contractAbi.Unpack(&approvalEvent, "Approval", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			approvalEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
			approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("Token Owner: %s\n", approvalEvent.TokenOwner.Hex())
			fmt.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
			fmt.Printf("Tokens: %s\n", approvalEvent.Tokens.String())
		}

		fmt.Printf("\n\n")
	}
}
*/
