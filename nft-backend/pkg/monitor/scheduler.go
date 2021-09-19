package monitor

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Tasks"
)

func StartScheduler(redis asynq.RedisConnOpt) {
	scheduler := asynq.NewScheduler(redis, nil)
	task := asynq.NewTask(Tasks.TypeBurnTokens, nil)
	_, err := scheduler.Register("@every 100s", task)
	if err != nil {
		log.Println(err)
	}
	err = scheduler.Run()
	if err != nil {
		log.Println(err)
	}
	err = scheduler.Run()
	if err != nil {
		log.Println(err)
	}
}
