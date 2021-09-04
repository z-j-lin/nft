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

//add to mint queue
//stores the address in a list
func (db *Database) Qmint(address, resourceID string) error {
	//create a key for the set
	Job := address + resourceID
	//push the job on the mintq list
	db.Client.LPush(context.TODO(), "MintQ", Job)
	return nil
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
	var State *objects.State
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
	return State, nil
}

func (db *Database) UpdateState(state *objects.State) error {
	mapp := make(map[string]string)
	//Hset only takes strings
	mapp["HighestFinalizedBlock"] = fmt.Sprint(state.HighestFinalizedBlock)
	mapp["HighestProcessedBlock"] = fmt.Sprint(state.HighestProcessedBlock)
	_, err := db.Client.HSet(context.TODO(), "State", mapp).Result()
	if err != nil {
		log.Println(err, "@DBUpdateState")
		return err
	}
	return err
}

func (db *Database) QPendingBlock(blocknum uint64) error {
	db.Client.LPush(context.TODO(), "PendingBlock", blocknum)
	return nil
}

func (db *Database) DQpendingBlock() (Blocknum uint64) {
	//BRPOP from the end of mint
	blocknum, err := db.Client.BRPop(context.TODO(), 1*time.Second, "PendingBlock").Result()
	if err != nil {
		log.Println("at DQpending", err)
	}
	Blocknum, err = strconv.ParseUint(blocknum[0], 10, 64)
	if err != nil {
		log.Panicf("failed to convert blocknum to uint: %v", err)
	}
	return
}

/*ownership store
each address gets a hmap
in the hmap each fied is a tokenID refrencing resourceID
*/
func (db *Database) StoreOwnership(resourceID, accountAddr, tokenID string, days2Live float64) {
	//
	data := make(map[string]string)
	data[tokenID] = resourceID
	_, err := db.Client.HSet(context.TODO(), accountAddr, data).Result()
	if err != nil {
		panic(err)
	}
	// add tokenID to the resourceID set
	err = db.Client.SAdd(context.TODO(), resourceID, tokenID).Err()
	if err != nil {
		panic(err)
	}
	//add token to collective token sorted set using expiration date to make the rank
	currentDay, err := db.Client.Get(context.TODO(), "day").Float64()
	if err != nil {
		panic(err)
	}
	element := &redis.Z{Score: currentDay + days2Live, Member: tokenID}
	//this is used to get the delete token array for burning tokens
	db.Client.ZAdd(context.TODO(), "Collective", element)
}

/*resource access store
make a set with resourceID as key, store tokenID in the set
on request of resource
*/
func (db *Database) RemoveToken(tokenID, OwnerAddr string) {
	//get resourceID from account map
	resourceID := db.Client.HMGet(context.TODO(), OwnerAddr, tokenID).String()
	//delete TokenID from resource set
	db.Client.SRem(context.TODO(), resourceID, tokenID)
	//delete tokenID from owner hash map
	db.Client.HDel(context.TODO(), OwnerAddr, tokenID)
}
