package monitor

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

type Validator struct {
	eth      *blockchain.Ethereum
	Blocknum *big.Int
	db       *redisDb.Database
	txhash   string
}

func NewValidator(eth *blockchain.Ethereum, blocknum *big.Int) error {
	rdb, err := redisDb.NewDBinstance()
	if err != nil {
		log.Fatal("error creating db instance, ", err)
		return err
	}
	val := &Validator{
		eth:      eth,
		Blocknum: blocknum,
		db:       rdb,
	}
	err = val.validateBlock()
	return err
}

func (v *Validator) validateBlock() error {
	//get the block
	block, err := v.eth.Client.BlockByNumber(context.TODO(), v.Blocknum)
	if err != nil {
		log.Println("error getting block at validator:", err)
		return err
	}
	//read receipt from block
	BlockTXs := block.Transactions()
	//iterate through the transactions looking for transactions to CAToken Contract
	contractAddr := v.eth.Contract.ContractAddress
	for _, transaction := range BlockTXs {
		to := *transaction.To()
		//looking for transactions we care about
		if to == contractAddr {
			Txhash := transaction.Hash()

			v.txhash = Txhash.Hex()
			TxReceipt, err := v.eth.Client.TransactionReceipt(context.TODO(), Txhash)
			if err != nil {
				log.Println("error getting transaction receipt at validator:", err)
				return err
			}
			//get the logs of the transaction
			TXlog := TxReceipt.Logs
			Log := TXlog[1]
			//decide what to do with the event
			v.EventHandler(Log)
		}
	}
	return nil
}

func (v *Validator) EventHandler(Log *types.Log) error {
	switch Log.Topics[0].Hex() {
	case v.eth.Contract.MintEvent:
		//extract the token recipient address
		RecipientAddr := Log.Topics[1].Hex()
		tokenID := Log.Topics[2].String()
		//get the resourceID from the contract
		resourceID, err := v.eth.Contract.GetResourceID(tokenID)
		if err != nil {
			log.Println("validator: error occured getting resourceID")
			return err
		}
		//add recipient Address and token id to registry
		err = v.db.StoreOwnership(resourceID, RecipientAddr, tokenID, 10)
		return err
	case v.eth.Contract.TransferEvent:
		from := Log.Topics[1].Hex()
		to := Log.Topics[2].String()
		tokenID := v.db.Client.Get(context.TODO(), v.txhash).String()
		fmt.Println("to Address", to)
		fmt.Println("from Address", from)
		fmt.Println("TokenID", tokenID)
		return nil
	default:
		return nil
	}

}
