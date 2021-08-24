package monitor

import (
	"context"
	"github.com/go-redis/redis/v8"

)


type Database struct{
	Client *redis.Client
}

ctx := context.TODO()
func NewDB(address string) (*Database, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: ""
		DB: 0,
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

//add to transaction queue
func (DB *Database) qTrans(resourceID, address string) error {
	DB.client.
}
//take off the transaction queue
func (DB *Database) DQTrans() error {

} 
//add to the verification queue 
func (DB *Database) qVerifie(tokenID, Resource, Address, Hash) error {

}
//take off the verification queue 
func (DB *Database) dqVerifie()