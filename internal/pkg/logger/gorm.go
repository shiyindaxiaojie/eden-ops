package logger

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	log *logrus.Logger
}

func NewGormLogger(log *logrus.Logger) *GormLogger {
	return &GormLogger{
		log: log,
	}
}

func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// 不记录 Info 级别的日志
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// 不记录 Warn 级别的日志
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)

	timestamp := time.Now().Format("2006/01/02 15:04:05.000")
	logMsg := fmt.Sprintf("%s %s:%d 数据库错误: %s", timestamp, file, line, fmt.Sprintf(msg, data...))
	l.log.Error(logMsg)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 获取调用信息
	_, file, line, ok := runtime.Caller(2) // 调整这个数字以获取正确的调用栈
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)

	// 构建日志格式：时间 文件:行号 [耗时ms] [rows:行数] SQL语句
	timestamp := time.Now().Format("2006/01/02 15:04:05.000")
	logMsg := fmt.Sprintf("%s %s:%d [%dms] [rows:%d] %s",
		timestamp,
		file,
		line,
		elapsed.Milliseconds(),
		rows,
		sql)

	if err != nil {
		l.log.WithField("error", err).Error(logMsg)
	} else {
		l.log.Info(logMsg)
	}
}
