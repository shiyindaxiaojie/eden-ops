package logger

import (
	"eden-ops/pkg/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	globalLogger *GlobalLogger
	oldLog       *logrus.Logger // 保持兼容性
	sqlEnabled   bool           // 是否启用SQL日志
	sqlDetailed  bool           // 是否显示详细SQL日志
)

// getProjectCaller 获取项目代码的调用者信息，跳过第三方库代码
func getProjectCaller() (string, int) {
	projectPath := "eden-ops"
	skipPaths := []string{
		"gorm.io",
		"github.com/sirupsen/logrus",
		"github.com/gin-gonic/gin",
		"runtime/",
		"asm_amd64.s",
		"proc.go",
		"pkg/logger/logger.go",
		"internal/pkg/logger/gorm.go",
		"finisher_api.go",
		"callbacks.go",
		"create.go",
		"update.go",
		"delete.go",
		"query.go",
		"raw.go",
		"transaction.go",
		"migrator.go",
	}

	for i := 1; i < 25; i++ { // 从调用栈的第1层开始查找，最多查找25层
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fileName := filepath.Base(file)

		// 首先检查是否为项目代码
		if !strings.Contains(file, projectPath) {
			continue
		}

		// 检查是否需要跳过的路径
		shouldSkip := false
		for _, skipPath := range skipPaths {
			if strings.Contains(file, skipPath) || strings.Contains(fileName, skipPath) {
				shouldSkip = true
				break
			}
		}

		// 如果不需要跳过，则返回
		if !shouldSkip {
			return fileName, line
		}
	}

	// 如果没找到项目代码，返回默认值
	_, file, line, _ := runtime.Caller(1)
	return filepath.Base(file), line
}

// GlobalLogger 全局日志器
type GlobalLogger struct {
	logger *logrus.Logger
}

// GlobalFormatter 全局日志格式化器
type GlobalFormatter struct{}

// Format 格式化日志
func (f *GlobalFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用信息
	_, file, line, _ := runtime.Caller(10) // 调整这个数字以获取正确的调用栈
	if file == "" {
		file = "unknown"
		line = 0
	}
	file = filepath.Base(file)

	// 获取当前时间（毫秒精度）
	timestamp := entry.Time.Format("2006/01/02 15:04:05.000")

	// 构建基础日志消息
	var logMessage string

	// 检查是否是 SQL 日志
	if sqlDuration, ok := entry.Data["sql_duration"].(time.Duration); ok {
		// 如果SQL日志被禁用，跳过输出
		if !sqlEnabled {
			return nil, nil
		}

		sql, _ := entry.Data["sql"].(string)
		rows, _ := entry.Data["rows"].(int64)
		isError, _ := entry.Data["sql_error"].(bool)

		if isError || sqlDetailed {
			// 错误时显示完整SQL，或者配置为详细模式
			logMessage = fmt.Sprintf("%s %s:%d [SQL] [%dms] [rows:%d] %s\n",
				timestamp, file, line, sqlDuration.Milliseconds(), rows, sql)
		} else {
			// 正常情况下限制SQL长度为100字符
			if len(sql) > 100 {
				sql = sql[:100] + "..."
			}
			logMessage = fmt.Sprintf("%s %s:%d [SQL] [%dms] [rows:%d] %s\n",
				timestamp, file, line, sqlDuration.Milliseconds(), rows, sql)
		}
	} else if apiDuration, ok := entry.Data["api_duration"].(time.Duration); ok {
		// API 接口日志
		method, _ := entry.Data["method"].(string)
		path, _ := entry.Data["path"].(string)
		statusCode, _ := entry.Data["status_code"].(int)
		clientIP, _ := entry.Data["client_ip"].(string)

		if clientIP == "::1" {
			clientIP = "127.0.0.1"
		}

		logMessage = fmt.Sprintf("%s %s:%d [API] [%dms] [%d] %s %s %s\n",
			timestamp, file, line, apiDuration.Milliseconds(), statusCode, clientIP, method, path)
	} else {
		// 普通日志
		logMessage = fmt.Sprintf("%s %s:%d %s\n",
			timestamp, file, line, entry.Message)
	}

	return []byte(logMessage), nil
}

// NewGlobalLogger 创建全局日志器
func NewGlobalLogger() *GlobalLogger {
	logger := logrus.New()
	logger.SetFormatter(&GlobalFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)

	return &GlobalLogger{
		logger: logger,
	}
}

// GlobalLogger 方法
func (l *GlobalLogger) Info(format string, args ...interface{}) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Infof(format, args...)
}

func (l *GlobalLogger) Error(format string, args ...interface{}) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Errorf(format, args...)
}

func (l *GlobalLogger) Warn(format string, args ...interface{}) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Warnf(format, args...)
}

func (l *GlobalLogger) Debug(format string, args ...interface{}) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Debugf(format, args...)
}

func (l *GlobalLogger) SQL(sql string, duration time.Duration, rows int64) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file":         file,
		"line":         line,
		"sql_duration": duration,
		"sql":          sql,
		"rows":         rows,
		"sql_error":    false,
	}).Info("SQL执行")
}

func (l *GlobalLogger) SQLWithError(sql string, duration time.Duration, rows int64, isError bool) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file":         file,
		"line":         line,
		"sql_duration": duration,
		"sql":          sql,
		"rows":         rows,
		"sql_error":    isError,
	}).Info("SQL执行")
}

func (l *GlobalLogger) API(method, path, clientIP string, statusCode int, duration time.Duration) {
	file, line := getProjectCaller()
	l.logger.WithFields(logrus.Fields{
		"file":         file,
		"line":         line,
		"api_duration": duration,
		"method":       method,
		"path":         path,
		"status_code":  statusCode,
		"client_ip":    clientIP,
	}).Info("API请求")
}

// 全局日志函数
func Info(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(format, args...)
	} else {
		// 降级到标准日志
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006/01/02 15:04:05.000")
		logMsg := fmt.Sprintf("%s %s:%d "+format, append([]interface{}{timestamp, filepath.Base(file), line}, args...)...)
		log.Print(logMsg)
	}
}

func Error(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(format, args...)
	} else {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006/01/02 15:04:05.000")
		logMsg := fmt.Sprintf("%s %s:%d ERROR "+format, append([]interface{}{timestamp, filepath.Base(file), line}, args...)...)
		log.Print(logMsg)
	}
}

func Warn(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(format, args...)
	} else {
		_, file, line, _ := runtime.Caller(1)
		timestamp := time.Now().Format("2006/01/02 15:04:05.000")
		logMsg := fmt.Sprintf("%s %s:%d WARN "+format, append([]interface{}{timestamp, filepath.Base(file), line}, args...)...)
		log.Print(logMsg)
	}
}

func Debug(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(format, args...)
	}
}

func SQL(sql string, duration time.Duration, rows int64) {
	if globalLogger != nil {
		globalLogger.SQL(sql, duration, rows)
	}
}

func SQLWithError(sql string, duration time.Duration, rows int64, isError bool) {
	if globalLogger != nil {
		globalLogger.SQLWithError(sql, duration, rows, isError)
	}
}

func API(method, path, clientIP string, statusCode int, duration time.Duration) {
	if globalLogger != nil {
		globalLogger.API(method, path, clientIP, statusCode, duration)
	}
}

// Init 初始化日志系统
func Init(cfg config.LogConfig) error {
	// 创建全局日志器
	globalLogger = NewGlobalLogger()

	// 设置SQL日志配置
	sqlEnabled = cfg.SQLEnabled
	sqlDetailed = cfg.SQLDetailed

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	globalLogger.logger.SetLevel(level)

	// 设置日志输出
	if cfg.Output == "file" && cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		globalLogger.logger.SetOutput(file)
	}

	// 保持兼容性
	oldLog = globalLogger.logger

	return nil
}

// GetLogger 获取日志实例（保持兼容性）
func GetLogger() *logrus.Logger {
	if oldLog != nil {
		return oldLog
	}
	if globalLogger != nil {
		return globalLogger.logger
	}
	return nil
}

// LogHTTPRequest 记录HTTP请求日志（保持兼容性）
func LogHTTPRequest(statusCode int, method, path string, latency time.Duration, clientIP string) {
	API(method, path, clientIP, statusCode, latency)
}

// CustomFormatter 自定义日志格式化器（保持兼容性）
type CustomFormatter struct{}

// Format 格式化日志
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return (&GlobalFormatter{}).Format(entry)
}
