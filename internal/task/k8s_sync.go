package task

import (
	"context"
	"eden-ops/internal/pkg/utils"
	"eden-ops/internal/service"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// K8sSyncTask Kubernetes同步任务
type K8sSyncTask struct {
	db         *gorm.DB
	service    service.K8sConfigService
	logger     *logrus.Logger
	cron       *cron.Cron
	jobEntries map[int64]cron.EntryID // 存储每个集群的任务ID
}

// NewK8sSyncTask 创建Kubernetes同步任务
func NewK8sSyncTask(db *gorm.DB, service service.K8sConfigService, logger *logrus.Logger) *K8sSyncTask {
	return &K8sSyncTask{
		db:         db,
		service:    service,
		logger:     logger,
		cron:       cron.New(cron.WithSeconds()), // 启用秒级支持
		jobEntries: make(map[int64]cron.EntryID),
	}
}

// Start 启动同步任务
func (t *K8sSyncTask) Start(ctx context.Context) error {
	t.logger.Infof("启动 Kubernetes 同步任务")

	// 启动定时检查任务，每30秒检查一次配置变化（6字段格式：秒 分 时 日 月 周）
	_, err := t.cron.AddFunc("*/30 * * * * *", func() {
		t.refreshSyncJobs()
	})

	if err != nil {
		return err
	}

	// 初始化同步任务
	t.refreshSyncJobs()

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

// RefreshJobs 立即刷新同步任务（供外部调用）
func (t *K8sSyncTask) RefreshJobs() {
	t.refreshSyncJobs()
}

// refreshSyncJobs 刷新同步任务
func (t *K8sSyncTask) refreshSyncJobs() {
	// 获取调用信息
	_, filePath, line, _ := runtime.Caller(0)
	file := filepath.Base(filePath)

	// 获取所有启用的Kubernetes配置
	configs, _, err := t.service.List(1, 1000, "", nil, nil, "")
	if err != nil {
		timestamp := time.Now().Format(utils.DateTimeMillisecond)
		t.logger.Errorf("%s %s:%d 获取Kubernetes配置列表失败: %v", timestamp, file, line, err)
		return
	}

	// 当前活跃的配置ID
	activeConfigIDs := make(map[int64]bool)

	// 为每个启用的集群创建或更新同步任务
	for _, config := range configs {
		configID := int64(config.ID)
		activeConfigIDs[configID] = true

		// 检查集群是否启用
		if config.Status != 1 {
			// 如果集群被禁用，移除其同步任务
			if entryID, exists := t.jobEntries[configID]; exists {
				t.cron.Remove(entryID)
				delete(t.jobEntries, configID)
				timestamp := time.Now().Format(utils.DateTimeMillisecond)
				t.logger.Infof("%s %s:%d 移除集群 %s 的同步任务（已禁用）", timestamp, file, line, config.Name)
			}
			continue
		}

		// 检查是否已经存在同步任务
		if _, exists := t.jobEntries[configID]; exists {
			continue // 任务已存在，跳过
		}

		// 创建新的同步任务
		syncInterval := config.SyncInterval
		if syncInterval < 30 {
			syncInterval = 30 // 最低30秒
		}

		// 构建cron表达式：每N秒执行一次（6字段格式：秒 分 时 日 月 周）
		cronExpr := fmt.Sprintf("*/%d * * * * *", syncInterval)

		entryID, err := t.cron.AddFunc(cronExpr, func() {
			t.syncSingleCluster(configID, config.Name)
		})

		if err != nil {
			timestamp := time.Now().Format(utils.DateTimeMillisecond)
			t.logger.Errorf("%s %s:%d 为集群 %s 创建同步任务失败: %v", timestamp, file, line, config.Name, err)
			continue
		}

		t.jobEntries[configID] = entryID
		timestamp := time.Now().Format(utils.DateTimeMillisecond)
		t.logger.Infof("%s %s:%d 为集群 %s 创建同步任务，间隔 %d 秒", timestamp, file, line, config.Name, syncInterval)
	}

	// 移除不再存在的集群的同步任务
	for configID, entryID := range t.jobEntries {
		if !activeConfigIDs[configID] {
			t.cron.Remove(entryID)
			delete(t.jobEntries, configID)
			timestamp := time.Now().Format(utils.DateTimeMillisecond)
			t.logger.Infof("%s %s:%d 移除集群 ID %d 的同步任务（配置已删除）", timestamp, file, line, configID)
		}
	}
}

// syncSingleCluster 同步单个集群
func (t *K8sSyncTask) syncSingleCluster(configID int64, clusterName string) {
	// 获取调用信息
	_, filePath, line, _ := runtime.Caller(0)
	file := filepath.Base(filePath)

	if err := t.service.SyncCluster(configID); err != nil {
		timestamp := time.Now().Format(utils.DateTimeMillisecond)
		t.logger.Errorf("%s %s:%d 同步集群 %s 失败: %v", timestamp, file, line, clusterName, err)
	} else {
		timestamp := time.Now().Format(utils.DateTimeMillisecond)
		t.logger.Infof("%s %s:%d 同步集群 %s 成功", timestamp, file, line, clusterName)
	}
}
