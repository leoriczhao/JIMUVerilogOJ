package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"verilog-oj/judge-service/internal/config"
	"verilog-oj/judge-service/internal/judge"
	"verilog-oj/judge-service/internal/queue"
)

func main() {
	// 加载配置
	cfg := config.LoadJudgeConfig()

	// 初始化判题器
	judger := judge.NewJudge(cfg.WorkDir)

	// 初始化消息队列
	rq := queue.NewRedisQueue(
		cfg.Queue.Host,
		cfg.Queue.Port,
		cfg.Queue.Password,
		cfg.Queue.DB,
		cfg.Queue.QueueName,
	)
	defer rq.Close()

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动判题服务
	log.Println("Starting judge service...")
	go startJudgeWorker(ctx, judger, rq)

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down judge service...")
	cancel()

	// 等待一段时间让正在进行的判题完成
	time.Sleep(5 * time.Second)
	log.Println("Judge service stopped")
}

// startJudgeWorker 启动判题工作器
func startJudgeWorker(ctx context.Context, judger *judge.Judge, rq *queue.RedisQueue) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 从队列获取判题请求
			request, err := rq.Pop(ctx)
			if err != nil {
				log.Printf("Failed to pop from queue: %v", err)
				continue
			}

			if request == nil {
				// 队列为空，继续轮询
				continue
			}

			log.Printf("Processing submission: %s", request.SubmissionID)

			// 执行判题
			result, err := judger.Judge(ctx, request)
			if err != nil {
				log.Printf("Judge failed for submission %s: %v", request.SubmissionID, err)
				continue
			}

			// 发布结果
			if err := rq.PublishResult(ctx, result); err != nil {
				log.Printf("Failed to publish result for submission %s: %v", request.SubmissionID, err)
				continue
			}

			log.Printf("Completed submission: %s, Status: %s, Score: %d", 
				result.SubmissionID, result.Status, result.Score)
		}
	}
} 