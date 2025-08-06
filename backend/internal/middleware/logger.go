package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// LogConfig 日志配置
type LogConfig struct {
	LogDir      string // 日志目录
	LogFile     string // 日志文件名
	MaxSize     int64  // 单个日志文件最大大小（字节）
	MaxAge      int    // 日志文件保留天数
	RotateDaily bool   // 是否按天轮转
}

// DefaultLogConfig 默认日志配置
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		LogDir:      "/root/logs",
		LogFile:     "api.log",
		MaxSize:     100 * 1024 * 1024, // 100MB
		MaxAge:      7,                  // 7天
		RotateDaily: true,
	}
}

// FileLogger 文件日志记录器
type FileLogger struct {
	config   *LogConfig
	currentFile *os.File
	currentDate string
}

// NewFileLogger 创建新的文件日志记录器
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

// openLogFile 打开日志文件
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
	return nil
}

// Write 写入日志
func (fl *FileLogger) Write(p []byte) (n int, err error) {
	// 检查是否需要轮转
	if fl.config.RotateDaily {
		now := time.Now()
		currentDate := now.Format("2006-01-02")
		if currentDate != fl.currentDate {
			if err := fl.openLogFile(); err != nil {
				return 0, err
			}
		}
	}

	// 检查文件大小
	if fl.config.MaxSize > 0 {
		if stat, err := fl.currentFile.Stat(); err == nil {
			if stat.Size() > fl.config.MaxSize {
				if err := fl.rotateBySize(); err != nil {
					return 0, err
				}
			}
		}
	}

	return fl.currentFile.Write(p)
}

// rotateBySize 按大小轮转日志文件
func (fl *FileLogger) rotateBySize() error {
	now := time.Now()
	timestamp := now.Format("20060102_150405")
	ext := filepath.Ext(fl.config.LogFile)
	name := fl.config.LogFile[:len(fl.config.LogFile)-len(ext)]
	newName := fmt.Sprintf("%s_%s%s", name, timestamp, ext)

	oldPath := filepath.Join(fl.config.LogDir, fl.config.LogFile)
	newPath := filepath.Join(fl.config.LogDir, newName)

	// 关闭当前文件
	fl.currentFile.Close()

	// 重命名文件
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("轮转日志文件失败: %v", err)
	}

	// 重新打开日志文件
	return fl.openLogFile()
}

// Close 关闭日志记录器
func (fl *FileLogger) Close() error {
	if fl.currentFile != nil {
		return fl.currentFile.Close()
	}
	return nil
}

// cleanOldLogs 清理过期日志文件
func (fl *FileLogger) cleanOldLogs() {
	if fl.config.MaxAge <= 0 {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -fl.config.MaxAge)

	filepath.Walk(fl.config.LogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			os.Remove(path)
		}

		return nil
	})
}

// FileLoggerMiddleware 文件日志中间件
func FileLoggerMiddleware(config *LogConfig) gin.HandlerFunc {
	// 强制输出调试信息，防止被编译器优化
	os.Stderr.WriteString("[DEBUG] 正在初始化文件日志\n")
	fileLogger, err := NewFileLogger(config)
	if err != nil {
		// 如果文件日志初始化失败，回退到标准输出
		os.Stderr.WriteString(fmt.Sprintf("[ERROR] 文件日志初始化失败: %v\n", err))
		return gin.Logger()
	}
	os.Stderr.WriteString("[DEBUG] 文件日志初始化成功\n")

	// 定期清理过期日志
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			fileLogger.cleanOldLogs()
		}
	}()

	// 只写入文件，不输出到标准输出
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Output: fileLogger,
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
	})
}

// APILogger API专用日志中间件
func APILogger() gin.HandlerFunc {
	return FileLoggerMiddleware(DefaultLogConfig())
}