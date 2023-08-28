package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendVerifyEmail = "task:send-verify-email"

type PayloadSendVerifyEmail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (publisher *RedisMessagePublisher) PublishTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := publisher.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Printf("type: %v, payload: %v, queue: %v, max_retry: %v -> enqueued task.",
		info.Type,
		payload,
		info.Queue,
		info.MaxRetry,
	)
	return nil
}
