package repository

import (
	"eden-ops/internal/model"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// K8sNodeHistoryRepository Node历史数据仓库接口
type K8sNodeHistoryRepository interface {
	// Node历史操作
	ArchiveNodesNotInList(configID int64, currentNodes []model.K8sNode, reason string) error
	GetNodeHistory(configID int64, page, pageSize int, startTime, endTime *time.Time) ([]model.K8sNodeHistory, int64, error)
	CleanupNodeHistory(beforeDate time.Time) error
	CountNodeHistory(configID int64) (int64, error)

	// 事务支持
	WithTx(tx *gorm.DB) K8sNodeHistoryRepository
	Transaction(fn func(K8sNodeHistoryRepository) error) error
}

// k8sNodeHistoryRepository Node历史数据仓库实现
type k8sNodeHistoryRepository struct {
	db *gorm.DB
}

// NewK8sNodeHistoryRepository 创建Node历史数据仓库
func NewK8sNodeHistoryRepository(db *gorm.DB) K8sNodeHistoryRepository {
	return &k8sNodeHistoryRepository{db: db}
}

// WithTx 使用事务
func (r *k8sNodeHistoryRepository) WithTx(tx *gorm.DB) K8sNodeHistoryRepository {
	return &k8sNodeHistoryRepository{db: tx}
}

// Transaction 执行事务
func (r *k8sNodeHistoryRepository) Transaction(fn func(K8sNodeHistoryRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTx(tx))
	})
}

// ArchiveNodesNotInList 归档不在当前列表中的Node
func (r *k8sNodeHistoryRepository) ArchiveNodesNotInList(configID int64, currentNodes []model.K8sNode, reason string) error {
	if len(currentNodes) == 0 {
		// 如果没有当前Node，归档所有Node
		sql := `INSERT INTO infra_k8s_node_history 
				(original_id, config_id, name, internal_ip, external_ip, hostname, os_image,
				 kernel_version, container_runtime, kubelet_version, kube_proxy_version,
				 cpu_capacity, memory_capacity, pods_capacity, cpu_allocatable, memory_allocatable,
				 pods_allocatable, cpu_usage, memory_usage, pods_usage, labels, annotations,
				 taints, conditions, status, ready, schedulable, created_at, updated_at,
				 deleted_at, archive_reason)
				SELECT id, config_id, name, internal_ip, external_ip, hostname, os_image,
					   kernel_version, container_runtime, kubelet_version, kube_proxy_version,
					   cpu_capacity, memory_capacity, pods_capacity, cpu_allocatable, memory_allocatable,
					   pods_allocatable, cpu_usage, memory_usage, pods_usage, labels, annotations,
					   taints, conditions, status, ready, schedulable, created_at, updated_at,
					   deleted_at, ?
				FROM infra_k8s_node 
				WHERE config_id = ? AND deleted_at IS NULL`
		return r.db.Exec(sql, reason, configID).Error
	}

	// 构建当前Node的名称列表
	var currentNames []string
	for _, node := range currentNodes {
		name := fmt.Sprintf("'%s'", strings.ReplaceAll(node.Name, "'", "''"))
		currentNames = append(currentNames, name)
	}

	// 归档不在当前列表中的Node
	sql := `INSERT INTO infra_k8s_node_history 
			(original_id, config_id, name, internal_ip, external_ip, hostname, os_image,
			 kernel_version, container_runtime, kubelet_version, kube_proxy_version,
			 cpu_capacity, memory_capacity, pods_capacity, cpu_allocatable, memory_allocatable,
			 pods_allocatable, cpu_usage, memory_usage, pods_usage, labels, annotations,
			 taints, conditions, status, ready, schedulable, created_at, updated_at,
			 deleted_at, archive_reason)
			SELECT id, config_id, name, internal_ip, external_ip, hostname, os_image,
				   kernel_version, container_runtime, kubelet_version, kube_proxy_version,
				   cpu_capacity, memory_capacity, pods_capacity, cpu_allocatable, memory_allocatable,
				   pods_allocatable, cpu_usage, memory_usage, pods_usage, labels, annotations,
				   taints, conditions, status, ready, schedulable, created_at, updated_at,
				   deleted_at, ?
			FROM infra_k8s_node 
			WHERE config_id = ? AND deleted_at IS NULL 
			AND name NOT IN (` + strings.Join(currentNames, ",") + `)`

	return r.db.Exec(sql, reason, configID).Error
}

// GetNodeHistory 获取Node历史记录
func (r *k8sNodeHistoryRepository) GetNodeHistory(configID int64, page, pageSize int, startTime, endTime *time.Time) ([]model.K8sNodeHistory, int64, error) {
	var histories []model.K8sNodeHistory
	var total int64

	query := r.db.Model(&model.K8sNodeHistory{}).Where("config_id = ?", configID)

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

// CleanupNodeHistory 清理Node历史记录
func (r *k8sNodeHistoryRepository) CleanupNodeHistory(beforeDate time.Time) error {
	return r.db.Where("archived_at < ?", beforeDate).Delete(&model.K8sNodeHistory{}).Error
}

// CountNodeHistory 统计Node历史记录数量
func (r *k8sNodeHistoryRepository) CountNodeHistory(configID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.K8sNodeHistory{}).Where("config_id = ?", configID).Count(&count).Error
	return count, err
}
