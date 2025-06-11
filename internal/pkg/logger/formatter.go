package logger

import (
	"eden-ops/internal/pkg/utils"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	TimestampFormat string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果是DEBUG级别，直接返回空
	if entry.Level == logrus.DebugLevel {
		return nil, nil
	}

	// 获取调用信息
	_, file, line, ok := runtime.Caller(6) // 调整这个数字以获取正确的调用栈
	if !ok {
		file = "???"
		line = 0
	}

	// 只保留文件名和行号
	file = filepath.Base(file)

	// 构建基本日志格式：时间 文件:行号
	timestamp := entry.Time.Format(utils.DateTimeMillisecond)
	msg := fmt.Sprintf("%s %s:%d", timestamp, file, line)

	// 如果有耗时信息，添加耗时
	if duration, ok := entry.Data["duration"]; ok {
		msg += fmt.Sprintf(" [%vms]", duration.(time.Duration).Milliseconds())
	}

	// 如果有影响的行数，添加行数信息
	if rows, ok := entry.Data["rows"]; ok {
		msg += fmt.Sprintf(" [rows:%v]", rows)
	}

	// 添加日志内容
	msg += fmt.Sprintf(" %s", entry.Message)

	// 如果有错误，添加错误信息
	if err, ok := entry.Data["error"]; ok {
		msg += fmt.Sprintf(": %v", err)
	}

	// 添加换行符
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	return []byte(msg), nil
}
