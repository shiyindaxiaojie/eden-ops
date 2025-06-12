package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// K8sWorkloadNamespaceRepository K8s工作负载命名空间仓库接口
type K8sWorkloadNamespaceRepository interface {
	GetByConfigID(configID int64) ([]*model.K8sWorkloadNamespace, error)
	CreateOrUpdate(namespace *model.K8sWorkloadNamespace) error
	DeleteByConfigID(configID int64) error
	UpdateWorkloadCount(configID int64, namespace string, count int) error
}

// k8sWorkloadNamespaceRepository K8s工作负载命名空间仓库实现
type k8sWorkloadNamespaceRepository struct {
	db *gorm.DB
}

// NewK8sWorkloadNamespaceRepository 创建K8s工作负载命名空间仓库
func NewK8sWorkloadNamespaceRepository(db *gorm.DB) K8sWorkloadNamespaceRepository {
	return &k8sWorkloadNamespaceRepository{db: db}
}

// GetByConfigID 根据配置ID获取命名空间列表
func (r *k8sWorkloadNamespaceRepository) GetByConfigID(configID int64) ([]*model.K8sWorkloadNamespace, error) {
	var namespaces []*model.K8sWorkloadNamespace
	err := r.db.Where("config_id = ?", configID).Find(&namespaces).Error
	return namespaces, err
}

// CreateOrUpdate 创建或更新命名空间
func (r *k8sWorkloadNamespaceRepository) CreateOrUpdate(namespace *model.K8sWorkloadNamespace) error {
	var existing model.K8sWorkloadNamespace
	err := r.db.Where("config_id = ? AND namespace = ?", namespace.ConfigID, namespace.Namespace).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		return r.db.Create(namespace).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有记录
		existing.WorkloadCount = namespace.WorkloadCount
		return r.db.Save(&existing).Error
	}
}

// DeleteByConfigID 删除指定配置的所有命名空间
func (r *k8sWorkloadNamespaceRepository) DeleteByConfigID(configID int64) error {
	return r.db.Where("config_id = ?", configID).Delete(&model.K8sWorkloadNamespace{}).Error
}

// UpdateWorkloadCount 更新命名空间的工作负载数量
func (r *k8sWorkloadNamespaceRepository) UpdateWorkloadCount(configID int64, namespace string, count int) error {
	return r.db.Model(&model.K8sWorkloadNamespace{}).
		Where("config_id = ? AND namespace = ?", configID, namespace).
		Update("workload_count", count).Error
}
