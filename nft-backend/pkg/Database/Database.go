package redisDb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Database struct {
	client *redis.Client
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
	Job := address + resourceID
	//push the job on the mintq list
	db.client.LPush(context.TODO(), "MintQ", Job)
	return nil
}

//take off the transaction queue
func (db *Database) DQmint() (account, resourceID string) {
	//BRPOP from the end of mint
	Job, err := db.client.BRPop(context.TODO(), 1*time.Second, "MintQ").Result()
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

//add to pending translist
//takes in the tx hash, recipient address, resource ID
func (db *Database) Qpending(txhash, address, resourceID string) error {
	tx := txhash + address + resourceID
	//push the job on the mintq list
	db.client.LPush(context.TODO(), "PendingTX", tx)
	return nil
}

func (db *Database) DQpending() (string, string, string) {
	//BRPOP from the end of mint
	Txdetails, err := db.client.BRPop(context.TODO(), 1*time.Second, "PendingTX").Result()
	if err != nil {
		panic(err)
	}
	txHash := Txdetails[1][:66]
	account := Txdetails[1][66:108]
	resourceID := Txdetails[1][108:]
	fmt.Println(Txdetails)
	return txHash, account, resourceID
}

func main() {
	rdb, err := NewDBinstance()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		Account, contentID := rdb.DQmint()

		fmt.Println("account:", Account, "contentID", contentID)
	}
}
