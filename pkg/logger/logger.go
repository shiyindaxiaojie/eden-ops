package logger

import (
	"eden-ops/pkg/config"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// CustomFormatter 自定义日志格式化器
type CustomFormatter struct {
}

// Format 实现logrus.Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取当前时间
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")

	// 构建日志消息
	// 对于非HTTP请求的日志，我们只显示时间和消息
	logMessage := fmt.Sprintf("%s | --- | ------------ | ------------- | ----- %q\n",
		timestamp,
		entry.Message,
	)

	// 如果有HTTP相关字段，则使用完整格式
	if statusCode, ok := entry.Data["status_code"].(int); ok {
		method, _ := entry.Data["method"].(string)
		path, _ := entry.Data["path"].(string)
		latency, _ := entry.Data["latency"].(time.Duration)
		clientIP, _ := entry.Data["client_ip"].(string)

		logMessage = fmt.Sprintf("%s | %3d | %12v | %15s | %-7s %q\n",
			timestamp,
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}

	return []byte(logMessage), nil
}

// Init 初始化日志
func Init(cfg config.LogConfig) error {
	log = logrus.New()

	// 使用自定义格式化器
	log.SetFormatter(&CustomFormatter{})

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	// 设置日志输出
	if cfg.Output == "file" && cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		log.SetOutput(file)
	}

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return log
}

// LogHTTPRequest 记录HTTP请求日志
func LogHTTPRequest(statusCode int, method, path string, latency time.Duration, clientIP string) {
	log.WithFields(logrus.Fields{
		"status_code": statusCode,
		"method":      method,
		"path":        path,
		"latency":     latency,
		"client_ip":   clientIP,
	}).Info("HTTP Request")
}
