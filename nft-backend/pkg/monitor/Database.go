package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Database struct {
	client *redis.Client
}

func NewDBinstance(address string) (*Database, error) {
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
		client: rdb,
	}, nil
}

func (db *Database) SetStringVal(key, val string) error {
	err := db.client.Set(context.TODO(), key, val, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
func (db *Database) getValue(key string) string {
	val, err := db.client.Get(context.TODO(), "key").Result()
	if err != nil {
		panic(err)
	}
	return val
}

//add to mint queue
//stores the address in a list
func (db *Database) Qmint(address, resourceID string) error {
	//create a key for the set
	key := address + resourceID
	//create a hash map with address as key holding account address and resourceID
	db.client.HSet(context.TODO(), key, "account", address)
	db.client.HSet(context.TODO(), key, "resourceID", resourceID)
	//push the key on the mintq list
	db.client.LPush(context.TODO(), "MintQ", key)
	return nil
}

//take off the transaction queue
func (db *Database) DQmint() {
	//BRPOP from the end of mint
	val, err := db.client.BRPop(context.TODO(), 10*time.Second, "MintQ").Result()
	if err != nil {
		fmt.Errorf("unable to pop of transaction list: %v", err)
	}
	fmt.Println(val)
}
