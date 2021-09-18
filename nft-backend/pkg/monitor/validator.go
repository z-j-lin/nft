package monitor

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	redisDb "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Database"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

const ZEROADDR string = "0x0000000000000000000000000000000000000000"

//const ZeroTokenID string = "0x0000000000000000000000000000000000000000000000000000000000000000"

type Validator struct {
	eth      *blockchain.Ethereum
	Blocknum *big.Int
	db       *redisDb.Database
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
		to := common.HexToAddress(ZEROADDR)
		if transaction.To() != nil {
			to = *transaction.To()
		}
		//looking for transactions we care about
		if to == contractAddr {
			//if its a transaction we care about loop through each of the logs
			//get the logs of the transaction
			Txhash := transaction.Hash()
			TxReceipt, err := v.eth.Client.TransactionReceipt(context.TODO(), Txhash)
			if err != nil {
				log.Println("Validator: error getting transaction receipt:", err)
				return err
			}
			TXlog := TxReceipt.Logs
			//loop through every log in the transaction
			for _, txlog := range TXlog {
				//look at the event of each log
				err := v.EventHandler(txlog)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// finds the blocks with a transfer event from the token Contract
func (v *Validator) EventHandler(Log *types.Log) error {
	var err error = nil
	//i only care about transfer events right now
	switch Log.Topics[0] {
	case v.eth.Contract.TransferEvent:
		log.Println("Transfer Event")
		from := Log.Topics[1].Hex()
		to := Log.Topics[2].Hex()
		//take the zeros out of the hash and convert to string
		tokenID := Log.Topics[3].Big().String()
		//trim the leading zero from the left side of the address string
		to = common.HexToAddress(to).Hex()
		from = common.HexToAddress(from).Hex()
		//if from is the 0 address
		//Mint event
		if from == ZEROADDR {
			log.Println("Mint")
			//add ownership to db
			err = v.db.StoreOwnership(to, tokenID, 10)

		} else if to == ZEROADDR { //burn event
			log.Println("Burn")
			//remove ownership
			err = v.db.DeleteOwnership(from, tokenID)

		} else { //transfering token
			log.Println("Transfering Ownership")
			err = v.db.TransferOwnership(from, to, tokenID)
		}
	default:
		err = nil
	}
	return err
}
