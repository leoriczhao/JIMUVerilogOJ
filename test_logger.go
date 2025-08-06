package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogConfig struct {
	LogDir      string
	LogFile     string
	MaxSize     int64
	MaxAge      int
	RotateDaily bool
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		LogDir:      "/root/logs",
		LogFile:     "api.log",
		MaxSize:     100 * 1024 * 1024, // 100MB
		MaxAge:      30,                 // 30天
		RotateDaily: true,
	}
}

type FileLogger struct {
	config      *LogConfig
	currentFile *os.File
	currentDate string
}

func NewFileLogger(config *LogConfig) (*FileLogger, error) {
	if config == nil {
		config = DefaultLogConfig()
	}

	// 创建日志目录
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	logger := &FileLogger{
		config: config,
	}

	// 打开日志文件
	if err := logger.openLogFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (fl *FileLogger) openLogFile() error {
	now := time.Now()
	var filename string

	if fl.config.RotateDaily {
		// 按日期轮转
		dateStr := now.Format("2006-01-02")
		ext := filepath.Ext(fl.config.LogFile)
		name := fl.config.LogFile[:len(fl.config.LogFile)-len(ext)]
		filename = fmt.Sprintf("%s_%s%s", name, dateStr, ext)
		fl.currentDate = dateStr
	} else {
		filename = fl.config.LogFile
	}

	filePath := filepath.Join(fl.config.LogDir, filename)
	fmt.Printf("尝试打开日志文件: %s\n", filePath)

	// 关闭当前文件
	if fl.currentFile != nil {
		fl.currentFile.Close()
	}

	// 打开新文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	fl.currentFile = file
	fmt.Printf("日志文件打开成功: %s\n", filePath)
	return nil
}

func (fl *FileLogger) Write(p []byte) (n int, err error) {
	return fl.currentFile.Write(p)
}

func main() {
	fmt.Println("测试日志初始化...")
	config := DefaultLogConfig()
	fmt.Printf("配置: LogDir=%s, LogFile=%s\n", config.LogDir, config.LogFile)
	
	logger, err := NewFileLogger(config)
	if err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		return
	}
	
	fmt.Println("日志初始化成功")
	
	// 写入测试日志
	testLog := "测试日志内容\n"
	n, err := logger.Write([]byte(testLog))
	if err != nil {
		fmt.Printf("写入日志失败: %v\n", err)
	} else {
		fmt.Printf("写入日志成功，字节数: %d\n", n)
	}
	
	logger.currentFile.Close()
}