package task

import (
	"context"
	"eden-ops/internal/service"
	"eden-ops/pkg/logger"
	"fmt"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// K8sSyncTask Kubernetes同步任务
type K8sSyncTask struct {
	db         *gorm.DB
	service    service.K8sConfigService
	cron       *cron.Cron
	jobEntries map[int64]cron.EntryID // 存储每个集群的任务ID
}

// NewK8sSyncTask 创建Kubernetes同步任务
func NewK8sSyncTask(db *gorm.DB, service service.K8sConfigService) *K8sSyncTask {
	return &K8sSyncTask{
		db:         db,
		service:    service,
		cron:       cron.New(cron.WithSeconds()), // 启用秒级支持
		jobEntries: make(map[int64]cron.EntryID),
	}
}

// Start 启动同步任务
func (t *K8sSyncTask) Start(ctx context.Context) error {
	logger.Info("启动 Kubernetes 同步任务")

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
		logger.Info("停止 Kubernetes 同步任务")
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
	// 获取所有启用的Kubernetes配置
	configs, _, err := t.service.List(1, 1000, "", nil, nil, "")
	if err != nil {
		logger.Error("获取Kubernetes配置列表失败: %v", err)
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
				logger.Info("移除集群 %s 的同步任务（已禁用）", config.Name)
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
			logger.Error("为集群 %s 创建同步任务失败: %v", config.Name, err)
			continue
		}

		t.jobEntries[configID] = entryID
		logger.Info("为集群 %s 创建同步任务，间隔 %d 秒", config.Name, syncInterval)
	}

	// 移除不再存在的集群的同步任务
	for configID, entryID := range t.jobEntries {
		if !activeConfigIDs[configID] {
			t.cron.Remove(entryID)
			delete(t.jobEntries, configID)
			logger.Info("移除集群 ID %d 的同步任务（配置已删除）", configID)
		}
	}
}

// syncSingleCluster 同步单个集群
func (t *K8sSyncTask) syncSingleCluster(configID int64, clusterName string) {
	if err := t.service.SyncCluster(configID); err != nil {
		logger.Error("同步集群 %s 失败: %v", clusterName, err)
	}
	// 成功信息已在SyncCluster方法中输出，这里不再重复记录
}
