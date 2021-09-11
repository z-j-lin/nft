package tasks

import (
	"encoding/json"
	"log"
	"time"

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

//run as go routine
// server just takes tasks off the queue
//numworkers determine max number of concurrent workers

type TaskClient struct {
	client *asynq.Client
}

func NewAsyncredisClient(redisAddr string) *TaskClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	return &TaskClient{
		client: client,
	}
}

// used in the api for token purchasing, creates a mint task
//only makes 1 mint task at a time
func (tc *TaskClient) QMintTask(accAddr, resourceID string) error {
	//create the task
	task, err := tc.newMintTokenTask(accAddr, resourceID)
	if err != nil {
		log.Println("failed to create Mint task", err)
		return err
	}
	// enques the mint task, retries every 10 min until succeeds or max retry count is hit
	info, err := tc.client.Enqueue(task, asynq.Queue("default"), asynq.Timeout(10*time.Minute), asynq.MaxRetry(3))
	if err != nil {
		log.Println("failed to queue the task")
		return err
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)
	return nil
}

//TODO: Implement QBURNTASK

func (tc *TaskClient) newMintTokenTask(AccAddr, resourceID string) (*asynq.Task, error) {
	Data, err := json.Marshal(MintToken{AccountAddress: AccAddr, ResourceID: resourceID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeMintToken, Data), nil
}

//might be able to make this a repeatable task with the asynq scheduler
func (tc *TaskClient) NewBurnTokenTask(tokenID string) (*asynq.Task, error) {
	Data, err := json.Marshal(BurnToken{tokenID: tokenID})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return asynq.NewTask(TypeBurnTokens, Data), nil
}
