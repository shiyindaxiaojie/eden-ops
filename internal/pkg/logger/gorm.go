package logger

import (
	"context"
	"eden-ops/pkg/logger"
	"time"

	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	log *logrus.Logger
}

func NewGormLogger(log *logrus.Logger) *GormLogger {
	return &GormLogger{
		log: log,
	}
}

func (l *GormLogger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// 不记录 Info 级别的日志
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// 不记录 Warn 级别的日志
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Error("数据库错误: "+msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		// 检查是否是"record not found"错误，这通常是正常的业务逻辑
		if err.Error() == "record not found" {
			// 对于"record not found"，只记录SQL日志，不记录错误
			logger.SQLWithError(sql, elapsed, rows, false)
		} else {
			// 其他错误才记录为错误日志
			logger.SQLWithError(sql, elapsed, rows, true)
			logger.Error("SQL执行错误: %v", err)
		}
	} else {
		// 正常情况下记录SQL日志
		logger.SQLWithError(sql, elapsed, rows, false)
	}
}
