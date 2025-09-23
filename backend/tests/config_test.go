package tests

import (
	"os"
	"testing"
)

// TestMain 在所有测试运行前后执行设置和清理工作
func TestMain(m *testing.M) {
	// 设置测试环境变量
	setupTestEnv()

	// 运行测试
	code := m.Run()

	// 清理测试环境
	cleanupTestEnv()

	// 退出
	os.Exit(code)
}

// setupTestEnv 设置测试环境
func setupTestEnv() {
	// 设置测试环境变量
	os.Setenv("GIN_MODE", "test")
	os.Setenv("LOG_LEVEL", "error")

	// 可以在这里设置测试数据库连接等
	// 例如：os.Setenv("DB_HOST", "localhost")
}

// cleanupTestEnv 清理测试环境
func cleanupTestEnv() {
	// 清理测试数据
	// 关闭数据库连接等
}

// TestConfig 测试配置结构
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
}

// GetTestConfig 获取测试配置
func GetTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("TEST_DB_USER", "test"),
		DBPassword: getEnvOrDefault("TEST_DB_PASSWORD", "test"),
		DBName:     getEnvOrDefault("TEST_DB_NAME", "test_db"),
		RedisHost:  getEnvOrDefault("TEST_REDIS_HOST", "localhost"),
		RedisPort:  getEnvOrDefault("TEST_REDIS_PORT", "6379"),
	}
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
