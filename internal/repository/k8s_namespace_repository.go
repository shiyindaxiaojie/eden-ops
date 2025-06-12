package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// K8sNamespaceRepository K8s命名空间仓库接口
type K8sNamespaceRepository interface {
	GetByConfigID(configID int64) ([]*model.K8sNamespace, error)
	CreateOrUpdate(namespace *model.K8sNamespace) error
	DeleteByConfigID(configID int64) error
	UpdateWorkloadCount(configID int64, namespace string, count int) error
}

// k8sNamespaceRepository K8s命名空间仓库实现
type k8sNamespaceRepository struct {
	db *gorm.DB
}

// NewK8sNamespaceRepository 创建K8s命名空间仓库
func NewK8sNamespaceRepository(db *gorm.DB) K8sNamespaceRepository {
	return &k8sNamespaceRepository{db: db}
}

// GetByConfigID 根据配置ID获取命名空间列表
func (r *k8sNamespaceRepository) GetByConfigID(configID int64) ([]*model.K8sNamespace, error) {
	var namespaces []*model.K8sNamespace
	err := r.db.Where("config_id = ?", configID).Find(&namespaces).Error
	return namespaces, err
}

// CreateOrUpdate 创建或更新命名空间
func (r *k8sNamespaceRepository) CreateOrUpdate(namespace *model.K8sNamespace) error {
	// 使用GORM的Clauses功能实现UPSERT操作，避免"record not found"错误
	return r.db.Where("config_id = ? AND namespace = ?", namespace.ConfigID, namespace.Namespace).
		Assign(model.K8sNamespace{WorkloadCount: namespace.WorkloadCount}).
		FirstOrCreate(namespace).Error
}

// DeleteByConfigID 删除指定配置的所有命名空间
func (r *k8sNamespaceRepository) DeleteByConfigID(configID int64) error {
	return r.db.Where("config_id = ?", configID).Delete(&model.K8sNamespace{}).Error
}

// UpdateWorkloadCount 更新命名空间的工作负载数量
func (r *k8sNamespaceRepository) UpdateWorkloadCount(configID int64, namespace string, count int) error {
	return r.db.Model(&model.K8sNamespace{}).
		Where("config_id = ? AND namespace = ?", configID, namespace).
		Update("workload_count", count).Error
}
