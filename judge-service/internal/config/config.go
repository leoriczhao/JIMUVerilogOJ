package config

import (
	"os"
	"strconv"
)

// JudgeConfig 判题服务配置
type JudgeConfig struct {
	WorkDir string      `yaml:"work_dir"`
	Queue   QueueConfig `yaml:"queue"`
}

// QueueConfig 消息队列配置
type QueueConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Password  string `yaml:"password"`
	DB        int    `yaml:"db"`
	QueueName string `yaml:"queue_name"`
}

// LoadJudgeConfig 加载判题服务配置
func LoadJudgeConfig() *JudgeConfig {
	return &JudgeConfig{
		WorkDir: getEnv("JUDGE_WORK_DIR", "/tmp/judge"),
		Queue: QueueConfig{
			Host:      getEnv("QUEUE_HOST", "localhost"),
			Port:      getEnvAsInt("QUEUE_PORT", 6379),
			Password:  getEnv("QUEUE_PASSWORD", ""),
			DB:        getEnvAsInt("QUEUE_DB", 0),
			QueueName: getEnv("QUEUE_NAME", "judge_queue"),
		},
	}
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
} 