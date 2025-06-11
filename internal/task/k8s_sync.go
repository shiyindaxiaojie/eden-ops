package task

import (
	"context"
	"eden-ops/internal/pkg/utils"
	"eden-ops/internal/service"
	"path/filepath"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// K8sSyncTask Kubernetes同步任务
type K8sSyncTask struct {
	db      *gorm.DB
	service service.K8sConfigService
	logger  *logrus.Logger
	cron    *cron.Cron
}

// NewK8sSyncTask 创建Kubernetes同步任务
func NewK8sSyncTask(db *gorm.DB, service service.K8sConfigService, logger *logrus.Logger) *K8sSyncTask {
	return &K8sSyncTask{
		db:      db,
		service: service,
		logger:  logger,
		cron:    cron.New(),
	}
}

// Start 启动同步任务
func (t *K8sSyncTask) Start(ctx context.Context) error {
	t.logger.Infof("启动 Kubernetes 同步任务")

	// 每5分钟同步一次
	_, err := t.cron.AddFunc("*/5 * * * *", func() {
		// 获取调用信息
		_, filePath, line, _ := runtime.Caller(0)
		file := filepath.Base(filePath)

		// 获取所有Kubernetes配置
		configs, _, err := t.service.List(1, 1000)
		if err != nil {
			timestamp := time.Now().Format(utils.DateTimeMillisecond)
			t.logger.Errorf("%s %s:%d 获取Kubernetes配置列表失败: %v", timestamp, file, line, err)
			return
		}

		// 同步每个集群
		for _, config := range configs {
			if err := t.service.SyncCluster(int64(config.ID)); err != nil {
				timestamp := time.Now().Format(utils.DateTimeMillisecond)
				t.logger.Errorf("%s %s:%d 同步集群 %s 失败: %v", timestamp, file, line, config.Name, err)
			} else {
				timestamp := time.Now().Format(utils.DateTimeMillisecond)
				t.logger.Infof("%s %s:%d 同步集群 %s 成功", timestamp, file, line, config.Name)
			}
		}
	})

	if err != nil {
		return err
	}

	t.cron.Start()

	// 监听上下文取消
	go func() {
		<-ctx.Done()
		t.Stop()

		// 获取调用信息
		_, filePath, line, _ := runtime.Caller(0)
		file := filepath.Base(filePath)

		timestamp := time.Now().Format(utils.DateTimeMillisecond)
		t.logger.Infof("%s %s:%d 停止 Kubernetes 同步任务", timestamp, file, line)
	}()

	return nil
}

// Stop 停止任务
func (t *K8sSyncTask) Stop() {
	if t.cron != nil {
		t.cron.Stop()
	}
}
