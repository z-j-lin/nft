package tasks

import (
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeMintToken  = "mintToken"
	TypeBurnTokens = "burn"
)

type MintToken struct {
	AccountAddress string
	ResourceID     string
	Nonce          uint64
}
type BurnToken struct {
	tokenID string
}

type AsyClient struct {
	client *asynq.Client
}

func NewServerClient(redisAddr string, numWorkers int) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 10},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeMintToken)
	mux.HandleFunc(TypeBurnTokens)
	err := srv.Run(mux)
}

func NewAsyncredisClient(redisAddr string) *AsyClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	return &AsyClient{
		client: client,
	}

}

func NewMintTokenTask(accAddr, resourceID string) (*asynq.Task, error) {
	Data, err := json.Marshal(MintToken{AccountAddress: accAddr, ResourceID: resourceID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeMintToken, Data), nil
}

func (ac *AsyClient) QMintTask(accAddr, resourceID string) error {
	//create the task
	task, err := NewMintTokenTask(accAddr, resourceID)
	if err != nil {
		log.Println("failed to create Mint task", err)
		return err
	}
	info, err := ac.client.Enqueue(task)
	if err != nil {
		log.Println("failed to queue the task")
		return err
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)
	return nil
}

func (ac *AsyClient) HandleMintTokenTask(t *asynq.Task) error {
	panic("unimplemented")
	//data struct stores data for the task
	var data MintToken
	err := json.Unmarshal(t.Payload(), &data)
	if err != nil {
		log.Println("failed to unmarshal task payload in Minttoken Handler")
		return err
	}
	//run send transaction function
	return nil
}

func (ac *AsyClient) NewBurnTokenTask(tokenID string) (*asynq.Task, error) {
	Data, err := json.Marshal(BurnToken{tokenID: tokenID})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return asynq.NewTask(TypeBurnTokens, Data), nil
}

func (ac *AsyClient) HandleBurnTokenTask(t *asynq.Task) error {
	panic("unimplemented")
	return nil
}
