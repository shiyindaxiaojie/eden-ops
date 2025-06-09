package logger

import (
	"context"
	"fmt"
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
	l.log.WithField("error", fmt.Sprintf(msg, data...)).Error("数据库错误")
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := logrus.Fields{
		"duration": elapsed,
	}

	if rows >= 0 {
		fields["rows"] = rows
	}

	if err != nil {
		fields["error"] = err
		l.log.WithFields(fields).Error(sql)
	} else {
		l.log.WithFields(fields).Info(sql)
	}
}
