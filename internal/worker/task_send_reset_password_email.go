package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendResetPasswordEmail = "task:send-password-reset-email"

type PayloadSendResetPasswordEmail struct {
	AppName string `json:"app_name"`
	Email   string `json:"email"`
	Token   string `json:"token"`
}

func (publisher *RedisMessagePublisher) PublishTaskSendResetPasswordEmail(ctx context.Context, payload *PayloadSendResetPasswordEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendResetPasswordEmail, jsonPayload, opts...)
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
