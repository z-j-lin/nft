package redisDb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	objects "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Objects"
)

type Database struct {
	Client *redis.Client
}

//creats a new redis client and returns a point to a Database object holding the client
func NewDBinstance() (*Database, error) {
	ctx := context.TODO()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	//checks if the database is live
	if err := rdb.Ping(ctx).Err(); err != nil {
		//if the database is not live return a empty pointer and err
		return nil, err
	}

	return &Database{
		Client: rdb,
	}, nil
}

func (db *Database) SetStringVal(key, val string) error {
	err := db.Client.Set(context.TODO(), key, val, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
func (db *Database) getValue(key string) string {
	val, err := db.Client.Get(context.TODO(), "key").Result()
	if err != nil {
		panic(err)
	}
	return val
}

//add to mint queue
//stores the address in a list
func (db *Database) Qmint(address, resourceID string) error {
	//create a key for the set
	Job := address + resourceID
	//push the job on the mintq list
	db.Client.LPush(context.TODO(), "MintQ", Job)
	return nil
}

//
func (db *Database) SetHighestFinalizedBlock(val string) error {
	err := db.Client.Set(context.TODO(), "HighestFinalizedBlock", val, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (db *Database) GetHighestFinalizedBlock(key string) string {
	val, err := db.Client.Get(context.TODO(), "HighestFinalizedBlock").Result()
	if err != nil {
		panic(err)
	}
	return val
}

//take off the transaction queue
func (db *Database) DQmint() (account, resourceID string) {
	//BRPOP from the end of mint
	Job, err := db.Client.BRPop(context.TODO(), 1*time.Second, "MintQ").Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err != redis.Nil {
		account := Job[1][:42]
		resourceID := Job[1][42:]
		return account, resourceID
	}
	return account, resourceID
}
func (db *Database) GetState() (*objects.State, error) {
	//fetch the state from redisDB
	StateSlice, err := db.Client.HMGet(context.TODO(), "State", "HighestFinalizedBlock", "HighestProcessedBlock").Result()
	if err != nil {
		panic(err)
	}
	var State objects.State
	State.HighestFinalizedBlock, err = strconv.ParseUint(StateSlice[0].(string), 10, 64)
	if err != nil {
		panic(err)
	}
	State.HighestProcessedBlock, err = strconv.ParseUint(StateSlice[1].(string), 10, 64)
	if err != nil {
		panic(err)
	}
	if StateSlice[1] == nil {
		logger := log.Default()
		logger.Println("State does not exist on disk")
		return nil, fmt.Errorf("state does not exist on disk")
	}
	if err != nil {
		log.Println(err, "@GetState")
	}

	fmt.Printf("stateSlice: %t\n", StateSlice[0])
	return &State, nil
}

func (db *Database) UpdateState(state objects.State) (int64, error) {
	mapp := make(map[string]string)
	//Hset only takes strings
	mapp["HighestFinalizedBlock"] = fmt.Sprint(state.HighestFinalizedBlock)
	mapp["HighestProcessedBlock"] = fmt.Sprint(state.HighestProcessedBlock)
	result, err := db.Client.HSet(context.TODO(), "State", mapp).Result()
	if err != nil {
		log.Println(err, "@DBUpdateState")
		return result, err
	}
	return result, err
}

func (db *Database) QPendingBlock(blocknum uint64) error {
	db.Client.LPush(context.TODO(), "PendingBlock", blocknum)
	return nil
}

//add to pending translist
//takes in the tx hash, recipient address, resource ID
func (db *Database) Qpending(txhash, address, resourceID string) error {
	tx := txhash + address + resourceID
	//push the job on the mintq list
	db.Client.LPush(context.TODO(), "PendingTX", tx)
	return nil
}

func (db *Database) DQpending() (txHash, account, resourceID string) {
	//BRPOP from the end of mint
	Txdetails, err := db.Client.BRPop(context.TODO(), 1*time.Second, "PendingTX").Result()
	if err != nil {
		log.Println("at DQpending", err)
	}
	txHash = Txdetails[1][:66]
	account = Txdetails[1][66:108]
	resourceID = Txdetails[1][108:]
	return txHash, account, resourceID
}
