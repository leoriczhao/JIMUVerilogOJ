package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"verilog-oj/judge-service/internal/judge"

	"github.com/go-redis/redis/v8"
)

// RedisQueue Redis消息队列实现
type RedisQueue struct {
	client    *redis.Client
	queueName string
}

// NewRedisQueue 创建Redis队列
func NewRedisQueue(host string, port int, password string, db int, queueName string) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &RedisQueue{
		client:    rdb,
		queueName: queueName,
	}
}

// Push 推送判题请求到队列
func (rq *RedisQueue) Push(ctx context.Context, request *judge.JudgeRequest) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	return rq.client.LPush(ctx, rq.queueName, data).Err()
}

// Pop 从队列弹出判题请求
func (rq *RedisQueue) Pop(ctx context.Context) (*judge.JudgeRequest, error) {
	result, err := rq.client.BRPop(ctx, 10*time.Second, rq.queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 队列为空
		}
		return nil, fmt.Errorf("failed to pop from queue: %v", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid result from queue")
	}

	var request judge.JudgeRequest
	if err := json.Unmarshal([]byte(result[1]), &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %v", err)
	}

	return &request, nil
}

// PublishResult 发布判题结果
func (rq *RedisQueue) PublishResult(ctx context.Context, result *judge.JudgeResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}

	resultChannel := fmt.Sprintf("judge_result_%s", result.SubmissionID)
	return rq.client.Publish(ctx, resultChannel, data).Err()
}

// SubscribeResults 订阅判题结果
func (rq *RedisQueue) SubscribeResults(ctx context.Context, submissionID string) (<-chan *judge.JudgeResult, error) {
	resultChannel := fmt.Sprintf("judge_result_%s", submissionID)
	pubsub := rq.client.Subscribe(ctx, resultChannel)

	resultChan := make(chan *judge.JudgeResult, 1)

	go func() {
		defer close(resultChan)
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					continue
				}

				var result judge.JudgeResult
				if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
					continue
				}

				select {
				case resultChan <- &result:
				case <-ctx.Done():
					return
				}
				return // 只接收一次结果
			}
		}
	}()

	return resultChan, nil
}

// Close 关闭连接
func (rq *RedisQueue) Close() error {
	return rq.client.Close()
}

// Health 健康检查
func (rq *RedisQueue) Health(ctx context.Context) error {
	return rq.client.Ping(ctx).Err()
} 