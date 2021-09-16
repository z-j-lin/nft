package Tasks

import (
	"encoding/json"
	"log"
	"math/big"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeMintToken        = "mintToken"
	TypeBurnTokens       = "burn"
	TypeBlockVerfication = "blockverification"
)

type MintToken struct {
	AccountAddress string
	ResourceID     string
	Nonce          uint64
}
type BurnToken struct {
	TokenIDs []*big.Int
}
type BlockV struct {
	Blocknum int64
}

type TaskClient struct {
	client *asynq.Client
}

func NewTaskClient(redisAddr string) *TaskClient {
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
	/* enques the mint task,
	retries every 10 min until succeeds
	or max retry count is hit*/
	info, err := tc.client.Enqueue(task, asynq.Queue("transactions"), asynq.Timeout(10*time.Minute), asynq.MaxRetry(3))
	if err != nil {
		log.Println("failed to queue the task")
		return err
	}
	log.Printf(" [*] Successfully enqueued Mint task: %+v", info)
	return nil
}

//task to be ran once 24 hours, should be used on a scheduler to periodicaly delete tokens
func (tc *TaskClient) QBurnTask(resourceIDs []*big.Int) error {
	//create the task
	task, err := tc.newBurnTokensTask(resourceIDs)
	if err != nil {
		log.Println("failed to create Mint task", err)
		return err
	}
	/* enques the burn tokens task,
	retries every 10 min until succeeds
	or max retry count is hit*/
	info, err := tc.client.Enqueue(
		task,
		asynq.Queue("default"),
		asynq.Timeout(10*time.Minute),
		asynq.MaxRetry(3),
	)
	if err != nil {
		log.Println("failed to queue burn task")
		return err
	}
	log.Printf(" [*] Successfully enqueued Burn task: %+v", info)
	return nil
}

//adds a verification task to a queue
func (tc *TaskClient) QVerificationTask(blocknum int64) error {
	//create the task
	task, err := tc.newBlockVerificationTask(blocknum)
	if err != nil {
		log.Println("failed to create Mint task", err)
		return err
	}
	/* enques the mint task,
	retries every 10 min until succeeds
	or max retry count is hit*/
	info, err := tc.client.Enqueue(
		task,
		asynq.Queue("validations"),
		asynq.Timeout(10*time.Minute),
		asynq.MaxRetry(3),
	)
	if err != nil {
		log.Println("failed to queue burn task")
		return err
	}
	log.Printf(" [*] Successfully enqueued verificcation task: %+v", info)
	return nil
}

func (tc *TaskClient) newMintTokenTask(AccAddr, resourceID string) (*asynq.Task, error) {
	Data, err := json.Marshal(MintToken{AccountAddress: AccAddr, ResourceID: resourceID})
	if err != nil {
		return nil, err
	}
	log.Printf(" [*] Successfully created a Mint task")
	return asynq.NewTask(TypeMintToken, Data), nil
}

//might be able to make this a repeatable task with the asynq scheduler
func (tc *TaskClient) newBurnTokensTask(tokenIDs []*big.Int) (*asynq.Task, error) {
	Data, err := json.Marshal(BurnToken{TokenIDs: tokenIDs})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf(" [*] Successfully created a Burn task")
	return asynq.NewTask(TypeBurnTokens, Data), nil
}

//validation task
func (tc *TaskClient) newBlockVerificationTask(blocknum int64) (*asynq.Task, error) {
	Data, err := json.Marshal(BlockV{Blocknum: blocknum})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf(" [*] Successfully created a verification task")
	return asynq.NewTask(TypeBlockVerfication, Data), nil
}
