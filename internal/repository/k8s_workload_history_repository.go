package repository

import (
	"eden-ops/internal/model"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// K8sWorkloadHistoryRepository Workload历史数据仓库接口
type K8sWorkloadHistoryRepository interface {
	// Workload历史操作
	ArchiveWorkloadsNotInList(configID int64, currentWorkloads []model.K8sWorkload, reason string) error
	GetWorkloadHistory(configID int64, page, pageSize int, startTime, endTime *time.Time) ([]model.K8sWorkloadHistory, int64, error)
	CleanupWorkloadHistory(beforeDate time.Time) error
	CountWorkloadHistory(configID int64) (int64, error)

	// 事务支持
	WithTx(tx *gorm.DB) K8sWorkloadHistoryRepository
	Transaction(fn func(K8sWorkloadHistoryRepository) error) error
}

// k8sWorkloadHistoryRepository Workload历史数据仓库实现
type k8sWorkloadHistoryRepository struct {
	db *gorm.DB
}

// NewK8sWorkloadHistoryRepository 创建Workload历史数据仓库
func NewK8sWorkloadHistoryRepository(db *gorm.DB) K8sWorkloadHistoryRepository {
	return &k8sWorkloadHistoryRepository{db: db}
}

// WithTx 使用事务
func (r *k8sWorkloadHistoryRepository) WithTx(tx *gorm.DB) K8sWorkloadHistoryRepository {
	return &k8sWorkloadHistoryRepository{db: tx}
}

// Transaction 执行事务
func (r *k8sWorkloadHistoryRepository) Transaction(fn func(K8sWorkloadHistoryRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTx(tx))
	})
}

// ArchiveWorkloadsNotInList 归档不在当前列表中的Workload
func (r *k8sWorkloadHistoryRepository) ArchiveWorkloadsNotInList(configID int64, currentWorkloads []model.K8sWorkload, reason string) error {
	if len(currentWorkloads) == 0 {
		// 如果没有当前Workload，归档所有Workload
		sql := `INSERT INTO infra_k8s_workload_history
				(original_id, config_id, name, namespace, kind, replicas, ready_replicas,
				 status, labels, selector, images, cpu_request, cpu_limit, memory_request,
				 memory_limit, created_at, updated_at, deleted_at, archive_reason)
				SELECT id, config_id, name, namespace, kind, replicas, ready_replicas,
					   status, labels, selector, images, cpu_request, cpu_limit, memory_request,
					   memory_limit, created_at, updated_at, deleted_at, ?
				FROM infra_k8s_workload
				WHERE config_id = ? AND deleted_at IS NULL`
		return r.db.Exec(sql, reason, configID).Error
	}

	// 构建当前Workload的唯一标识列表 (name-namespace-kind)
	var currentKeys []string
	for _, workload := range currentWorkloads {
		key := fmt.Sprintf("'%s-%s-%s'",
			strings.ReplaceAll(workload.Name, "'", "''"),
			strings.ReplaceAll(workload.Namespace, "'", "''"),
			strings.ReplaceAll(workload.Kind, "'", "''"))
		currentKeys = append(currentKeys, key)
	}

	// 归档不在当前列表中的Workload
	sql := `INSERT INTO infra_k8s_workload_history
			(original_id, config_id, name, namespace, kind, replicas, ready_replicas,
			 status, labels, selector, images, cpu_request, cpu_limit, memory_request,
			 memory_limit, created_at, updated_at, deleted_at, archive_reason)
			SELECT id, config_id, name, namespace, kind, replicas, ready_replicas,
				   status, labels, selector, images, cpu_request, cpu_limit, memory_request,
				   memory_limit, created_at, updated_at, deleted_at, ?
			FROM infra_k8s_workload
			WHERE config_id = ? AND deleted_at IS NULL
			AND CONCAT(name, '-', namespace, '-', kind) NOT IN (` + strings.Join(currentKeys, ",") + `)`

	return r.db.Exec(sql, reason, configID).Error
}

// GetWorkloadHistory 获取Workload历史记录
func (r *k8sWorkloadHistoryRepository) GetWorkloadHistory(configID int64, page, pageSize int, startTime, endTime *time.Time) ([]model.K8sWorkloadHistory, int64, error) {
	var histories []model.K8sWorkloadHistory
	var total int64

	query := r.db.Model(&model.K8sWorkloadHistory{}).Where("config_id = ?", configID)

	if startTime != nil {
		query = query.Where("archived_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("archived_at <= ?", *endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("archived_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// CleanupWorkloadHistory 清理Workload历史记录
func (r *k8sWorkloadHistoryRepository) CleanupWorkloadHistory(beforeDate time.Time) error {
	return r.db.Where("archived_at < ?", beforeDate).Delete(&model.K8sWorkloadHistory{}).Error
}

// CountWorkloadHistory 统计Workload历史记录数量
func (r *k8sWorkloadHistoryRepository) CountWorkloadHistory(configID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.K8sWorkloadHistory{}).Where("config_id = ?", configID).Count(&count).Error
	return count, err
}
