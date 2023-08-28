package worker

import (
	"github.com/hibiken/asynq"
)

type RedisMessagePublisher struct {
	client *asynq.Client
}

func NewRedisMessagePublisher(redisOptions asynq.RedisClientOpt) MessagePublisher {
	client := asynq.NewClient(redisOptions)
	return &RedisMessagePublisher{
		client: client,
	}
}
