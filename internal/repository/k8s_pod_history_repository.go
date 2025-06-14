package repository

import (
	"eden-ops/internal/model"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// K8sPodHistoryRepository Pod历史数据仓库接口
type K8sPodHistoryRepository interface {
	// Pod历史操作
	ArchivePodsNotInList(configID int64, currentPods []model.K8sPod, reason string) error
	GetPodHistory(configID int64, page, pageSize int, startTime, endTime *time.Time) ([]model.K8sPodHistory, int64, error)
	CleanupPodHistory(beforeDate time.Time) error
	CountPodHistory(configID int64) (int64, error)

	// 事务支持
	WithTx(tx *gorm.DB) K8sPodHistoryRepository
	Transaction(fn func(K8sPodHistoryRepository) error) error
}

// k8sPodHistoryRepository Pod历史数据仓库实现
type k8sPodHistoryRepository struct {
	db *gorm.DB
}

// NewK8sPodHistoryRepository 创建Pod历史数据仓库
func NewK8sPodHistoryRepository(db *gorm.DB) K8sPodHistoryRepository {
	return &k8sPodHistoryRepository{db: db}
}

// WithTx 使用事务
func (r *k8sPodHistoryRepository) WithTx(tx *gorm.DB) K8sPodHistoryRepository {
	return &k8sPodHistoryRepository{db: tx}
}

// Transaction 执行事务
func (r *k8sPodHistoryRepository) Transaction(fn func(K8sPodHistoryRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTx(tx))
	})
}

// ArchivePodsNotInList 归档不在当前列表中的Pod
func (r *k8sPodHistoryRepository) ArchivePodsNotInList(configID int64, currentPods []model.K8sPod, reason string) error {
	if len(currentPods) == 0 {
		// 如果没有当前Pod，归档所有Pod
		sql := `INSERT INTO infra_k8s_pod_history 
				(original_id, config_id, workload_id, name, namespace, workload_name, workload_kind,
				 status, phase, node_name, pod_ip, host_ip, instance_ip, cpu_request, cpu_limit,
				 memory_request, memory_limit, restart_count, start_time, created_at, updated_at,
				 deleted_at, archive_reason)
				SELECT id, config_id, workload_id, name, namespace, workload_name, workload_kind,
					   status, phase, node_name, pod_ip, host_ip, instance_ip, cpu_request, cpu_limit,
					   memory_request, memory_limit, restart_count, start_time, created_at, updated_at,
					   deleted_at, ?
				FROM infra_k8s_pod 
				WHERE config_id = ? AND deleted_at IS NULL`
		return r.db.Exec(sql, reason, configID).Error
	}

	// 构建当前Pod的唯一标识列表 (name-namespace)
	var currentKeys []string
	for _, pod := range currentPods {
		key := fmt.Sprintf("'%s-%s'", strings.ReplaceAll(pod.Name, "'", "''"), strings.ReplaceAll(pod.Namespace, "'", "''"))
		currentKeys = append(currentKeys, key)
	}

	// 归档不在当前列表中的Pod
	sql := `INSERT INTO infra_k8s_pod_history 
			(original_id, config_id, workload_id, name, namespace, workload_name, workload_kind,
			 status, phase, node_name, pod_ip, host_ip, instance_ip, cpu_request, cpu_limit,
			 memory_request, memory_limit, restart_count, start_time, created_at, updated_at,
			 deleted_at, archive_reason)
			SELECT id, config_id, workload_id, name, namespace, workload_name, workload_kind,
				   status, phase, node_name, pod_ip, host_ip, instance_ip, cpu_request, cpu_limit,
				   memory_request, memory_limit, restart_count, start_time, created_at, updated_at,
				   deleted_at, ?
			FROM infra_k8s_pod 
			WHERE config_id = ? AND deleted_at IS NULL 
			AND CONCAT(name, '-', namespace) NOT IN (` + strings.Join(currentKeys, ",") + `)`

	return r.db.Exec(sql, reason, configID).Error
}

// GetPodHistory 获取Pod历史记录
func (r *k8sPodHistoryRepository) GetPodHistory(configID int64, page, pageSize int, startTime, endTime *time.Time) ([]model.K8sPodHistory, int64, error) {
	var histories []model.K8sPodHistory
	var total int64

	query := r.db.Model(&model.K8sPodHistory{}).Where("config_id = ?", configID)

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

// CleanupPodHistory 清理Pod历史记录
func (r *k8sPodHistoryRepository) CleanupPodHistory(beforeDate time.Time) error {
	return r.db.Where("archived_at < ?", beforeDate).Delete(&model.K8sPodHistory{}).Error
}

// CountPodHistory 统计Pod历史记录数量
func (r *k8sPodHistoryRepository) CountPodHistory(configID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.K8sPodHistory{}).Where("config_id = ?", configID).Count(&count).Error
	return count, err
}
