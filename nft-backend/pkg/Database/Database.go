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

func (db *Database) GetAccTokens(accountAddr string) ([]string, error) {
	AccTokens, err := db.Client.SMembers(context.TODO(), accountAddr).Result()
	if err != nil {
		return AccTokens, err
	}
	return AccTokens, nil
}

//this function is used when a burn happens
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

func (db *Database) GetState() (objects.State, error) {
	//fetch the state from redisDB
	StateSlice, err := db.Client.HMGet(context.TODO(), "State", "HighestFinalizedBlock", "HighestProcessedBlock").Result()
	//if state doesnt exist return values will be all nil for state and err
	if err != nil {
		panic(err)
	}
	//if the state variables are not empty on the redis server
	var State objects.State
	if StateSlice[1] != nil {
		State.HighestProcessedBlock, err = strconv.ParseUint(StateSlice[1].(string), 10, 64)
		if err != nil {
			panic(err)
		}
		if StateSlice[0] != nil {
			State.HighestFinalizedBlock, err = strconv.ParseUint(StateSlice[0].(string), 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return State, nil
	} else {
		log.Println("State does not exist on disk")
		return State, fmt.Errorf("state does not exist on disk")
	}
}

////////////////////////////////////////////////////////////////////////////
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
func (db *Database) AddItem(contentID string) error {
	err := db.Client.SAdd(context.TODO(), "store", contentID).Err()
	return err
}

func (db *Database) GetStore() ([]string, error) {
	store, err := db.Client.SMembers(context.TODO(), "store").Result()
	if err != nil {
		return store, err
	}
	return store, nil
}

/*ownership store
each address gets a hmap
in the hmap each field is a tokenID refrencing resourceID
*/
func (db *Database) StoreOwnership(resourceID, accountAddr, tokenID string, days2Live float64) error {
	//each token has a hash map with the tokenID as the key
	//holds resource ID owners address
	//the days2live will be stored sorted set
	TokenHashData := make(map[string]string)
	TokenHashData["Resource"] = resourceID
	TokenHashData["Owner"] = accountAddr
	//create a tokenID map with tokenID as key, with fields resourceID, AccountOwners
	//used for serving content ownership array to the client provided the tokenID array
	//this is also used to verify access rights
	err := db.Client.HSet(context.TODO(), tokenID, TokenHashData).Err()

	if err != nil {
		log.Printf("db: failed to create token hash: %v\n", err)
		return err
	}
	//stores tokenID into owners set. this is done so we dont have to iterate through all the takenID hashmaps to find tokens
	//this is used for loading the owned page on the client side
	err = db.Client.SAdd(context.TODO(), accountAddr, tokenID).Err()
	if err != nil {
		log.Printf("db: failed to store tokenID into owners set: %v\n", err)
		return err
	}
	//add token to collective token sorted set ranked by expiration date
	days := int(days2Live)
	lifetime := time.Now().AddDate(0, 0, days)

	element := &redis.Z{Score: float64(lifetime.Unix()), Member: tokenID}
	//this is used to get the delete token array for burning tokens
	//the collection is token ID ranked with days2live
	err = db.Client.ZAdd(context.TODO(), "Collective", element).Err()
	if err != nil {
		log.Printf("db: failed to add the tokenID to the collective set: %v\n", err)
		return err
	}
	return nil
}
func (db *Database) DeleteOwnership(accountAddr, tokenID string, days2Live float64) error {
	err := db.Client.HDel(context.TODO(), tokenID, "Resource", "Owner").Err()
	if err != nil {
		return err
	}
	err = db.Client.SRem(context.TODO(), accountAddr, tokenID).Err()
	if err != nil {
		return err
	}
	return nil
}
