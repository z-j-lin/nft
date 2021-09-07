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

////////////////////////////////////////////////////////////////////////////

//add to mint queue
//stores the address in a list
func (db *Database) Qmint(address, resourceID string) error {
	//create a key for the set
	Job := address + resourceID
	//push the job on the mintq list
	numElem, err := db.Client.LPush(context.TODO(), "MintQ", Job).Result()
	if err != nil {
		return err
	}
	if numElem != 1 {
		panic("should not happen")
	}
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

//this function is used when a burn a burn happened
func (db *Database) RemoveToken(tokenID string) {
	//get Account from token map
	OwnerAddress := db.Client.HMGet(context.TODO(), tokenID, "Account").String()
	//delete tokenID from owner set
	err := db.Client.SRem(context.TODO(), OwnerAddress, tokenID).Err()
	if err != nil {
		log.Panicf("failed to delete token from account owner set: %v", err)
	}
	//delete token hash map
	err = db.Client.HDel(context.TODO(), tokenID).Err()
	if err != nil {
		log.Panicf("failed to delete token from hashmap: %v", err)
	}
}

/////////////////////////////////////////////////////////////////////////////////

//this function can add or just update blockmon state
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
	//need a condition for 0
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
	//each token has a hash map with the tokenID as the key
	//holds resource ID owners address
	//the days2live will be stored sorted set
	TokenHashData := make(map[string]string)
	TokenHashData["resource"] = resourceID
	TokenHashData["Owner"] = accountAddr
	//create a tokenID map with tokenID as key, with fields resourceID, AccountOwners
	//used for serving content ownership array to the client provided the tokenID array
	//this is also used to verify access rights
	err := db.Client.HSet(context.TODO(), tokenID, TokenHashData).Err()

	if err != nil {
		log.Panicf("failed to create token hash: %v", err)
	}
	//stores tokenID into owners set. this is done so we dont have to iterate through all the takenID hashmaps to find tokens
	//this is used for loading the owned page on the client side
	err = db.Client.SAdd(context.TODO(), accountAddr, tokenID).Err()
	if err != nil {
		log.Panicf("failed to store tokenID into owners set: %v", err)
	}
	//add token to collective token sorted set ranked by expiration date
	currentDay, err := db.Client.Get(context.TODO(), "day").Float64()
	if err != nil {
		log.Panicf("failed to get the current day from redisDB: %v", err)
	}
	element := &redis.Z{Score: currentDay + days2Live, Member: tokenID}
	//this is used to get the delete token array for burning tokens
	//the collection is token ID ranked with days2live
	err = db.Client.ZAdd(context.TODO(), "Collective", element).Err()
	if err != nil {
		log.Panicf("failed to add the tokenID to the collective set: %v", err)
	}
}
