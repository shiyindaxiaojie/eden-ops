package logger

import (
	"eden-ops/pkg/config"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// CustomFormatter 自定义日志格式化器
type CustomFormatter struct {
}

// Format 实现logrus.Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用信息
	_, file, line, ok := runtime.Caller(6) // 调整这个数字以获取正确的调用栈
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)

	// 获取当前时间
	timestamp := entry.Time.Format("2006/01/02 15:04:05.000")

	// 构建日志消息
	var logMessage string

	// 如果有HTTP相关字段，则使用HTTP请求格式
	if statusCode, ok := entry.Data["status_code"].(int); ok {
		method, _ := entry.Data["method"].(string)
		path, _ := entry.Data["path"].(string)
		latency, _ := entry.Data["latency"].(time.Duration)
		clientIP, _ := entry.Data["client_ip"].(string)

		// 如果是本地请求，将::1转换为127.0.0.1
		if clientIP == "::1" {
			clientIP = "127.0.0.1"
		}

		logMessage = fmt.Sprintf("%s | %3d | %dms | %15s | %-7s %q\n",
			timestamp,
			statusCode,
			latency.Milliseconds(),
			clientIP,
			method,
			path,
		)
	} else {
		// 对于普通日志，使用文件:行号格式
		logMessage = fmt.Sprintf("%s %s:%d %s\n",
			timestamp,
			file,
			line,
			entry.Message,
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
