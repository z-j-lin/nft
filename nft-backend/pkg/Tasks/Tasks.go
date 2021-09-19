package Tasks

import (
	"encoding/json"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/blockchain"
)

const (
	TypeMintToken        = "mintToken"
	TypeBurnTokens       = "burn"
	TypeBlockVerfication = "blockverification"
)

type MintToken struct {
	AccountAddress string
	ResourceID     string
	Nonce          *big.Int
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

func NewTaskClient(eth *blockchain.Ethereum, redisAddr string) *TaskClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &TaskClient{
		client: client,
	}
}

// used in the api for token purchasing, creates a mint task
//only makes 1 mint task at a time
func (tc *TaskClient) QMintTask(accAddr, resourceID string, Nonce *big.Int) error {
	//create the task

	task, err := tc.newMintTokenTask(accAddr, resourceID, Nonce)
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
func (tc *TaskClient) QBurnTask() error {
	//create the task
	task, err := tc.newBurnTokensTask()
	if err != nil {
		log.Println("failed to create Mint task", err)
		return err
	}
	/* enques the burn tokens task,
	retries every 10 min until succeeds
	or max retry count is hit*/
	info, err := tc.client.Enqueue(
		task,
		asynq.Queue("burn"),
		asynq.Timeout(10*time.Minute),
		asynq.MaxRetry(3),
		//run this every 24 hour
		asynq.ProcessAt(time.Now().Add(10*time.Minute)),
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
		log.Println("failed to create verification task", err)
		return err
	}
	/* enques the mint task,
	retries until succeeds
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

func (tc *TaskClient) newMintTokenTask(AccAddr, resourceID string, nonce *big.Int) (*asynq.Task, error) {
	Data, err := json.Marshal(MintToken{AccountAddress: AccAddr, ResourceID: resourceID, Nonce: nonce})
	if err != nil {
		return nil, err
	}
	log.Printf(" [*] Successfully created a Mint task")
	return asynq.NewTask(TypeMintToken, Data), nil
}

//might be able to make this a repeatable task with the asynq scheduler
func (tc *TaskClient) newBurnTokensTask() (*asynq.Task, error) {
	log.Println(" [*] Successfully created a Burn task")
	return asynq.NewTask(TypeBurnTokens, nil), nil
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

type NonceMan struct {
	sync.Mutex
	nonce *big.Int
}

func NewNonceManager(eth *blockchain.Ethereum) *NonceMan {
	//get next nonce from contract
	nonce, err := eth.Contract.GetInitNonce()
	log.Println("TXQmon: started new nonce manager with nonce", nonce)
	if err != nil {
		log.Panic(err)
	}
	return &NonceMan{
		nonce: nonce,
	}
}
func (nm *NonceMan) GetnonceWithLock() *big.Int {
	nm.Lock()
	defer nm.Unlock()
	nonce := nm.nonce
	nm.nonce = nonce.Add(nonce, big.NewInt(1))
	return nonce
}
