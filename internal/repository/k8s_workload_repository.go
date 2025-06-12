package repository

import (
	"eden-ops/internal/model"
	"gorm.io/gorm"
)

// K8sWorkloadRepository Kubernetes工作负载仓库接口
type K8sWorkloadRepository interface {
	Create(workload *model.K8sWorkload) error
	Update(workload *model.K8sWorkload) error
	Delete(id int64) error
	Get(id int64) (*model.K8sWorkload, error)
	List(configID int64, page, pageSize int) (int64, []model.K8sWorkload, error)
	ListWithFilter(page, pageSize int, name, namespace, workloadType string, configId *int64) (int64, []model.K8sWorkload, error)
	ListByConfigID(configID int64) ([]model.K8sWorkload, error)
	DeleteByConfigID(configID int64) error
	BatchCreate(workloads []model.K8sWorkload) error
	BatchUpdate(workloads []model.K8sWorkload) error
	CountByConfigID(configID int64) (int64, error)
}

// k8sWorkloadRepository Kubernetes工作负载仓库实现
type k8sWorkloadRepository struct {
	db *gorm.DB
}

// NewK8sWorkloadRepository 创建Kubernetes工作负载仓库实例
func NewK8sWorkloadRepository(db *gorm.DB) K8sWorkloadRepository {
	return &k8sWorkloadRepository{db: db}
}

// Create 创建工作负载
func (r *k8sWorkloadRepository) Create(workload *model.K8sWorkload) error {
	return r.db.Create(workload).Error
}

// Update 更新工作负载
func (r *k8sWorkloadRepository) Update(workload *model.K8sWorkload) error {
	return r.db.Save(workload).Error
}

// Delete 删除工作负载
func (r *k8sWorkloadRepository) Delete(id int64) error {
	return r.db.Delete(&model.K8sWorkload{}, id).Error
}

// Get 获取工作负载
func (r *k8sWorkloadRepository) Get(id int64) (*model.K8sWorkload, error) {
	var workload model.K8sWorkload
	err := r.db.First(&workload, id).Error
	if err != nil {
		return nil, err
	}
	return &workload, nil
}

// List 获取工作负载列表
func (r *k8sWorkloadRepository) List(configID int64, page, pageSize int) (int64, []model.K8sWorkload, error) {
	var workloads []model.K8sWorkload
	var total int64

	query := r.db.Model(&model.K8sWorkload{}).Where("config_id = ?", configID)
	
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&workloads).Error; err != nil {
		return 0, nil, err
	}

	return total, workloads, nil
}

// ListWithFilter 获取工作负载列表（支持筛选）
func (r *k8sWorkloadRepository) ListWithFilter(page, pageSize int, name, namespace, workloadType string, configId *int64) (int64, []model.K8sWorkload, error) {
	var workloads []model.K8sWorkload
	var total int64

	query := r.db.Model(&model.K8sWorkload{})

	// 添加筛选条件
	if configId != nil {
		query = query.Where("config_id = ?", *configId)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if namespace != "" {
		query = query.Where("namespace = ?", namespace)
	}
	if workloadType != "" {
		query = query.Where("kind = ?", workloadType)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&workloads).Error; err != nil {
		return 0, nil, err
	}

	return total, workloads, nil
}

// ListByConfigID 根据配置ID获取所有工作负载
func (r *k8sWorkloadRepository) ListByConfigID(configID int64) ([]model.K8sWorkload, error) {
	var workloads []model.K8sWorkload
	err := r.db.Where("config_id = ?", configID).Find(&workloads).Error
	return workloads, err
}

// DeleteByConfigID 根据配置ID删除所有工作负载
func (r *k8sWorkloadRepository) DeleteByConfigID(configID int64) error {
	return r.db.Where("config_id = ?", configID).Delete(&model.K8sWorkload{}).Error
}

// BatchCreate 批量创建工作负载
func (r *k8sWorkloadRepository) BatchCreate(workloads []model.K8sWorkload) error {
	if len(workloads) == 0 {
		return nil
	}
	return r.db.Create(&workloads).Error
}

// CountByConfigID 根据配置ID统计工作负载数量
func (r *k8sWorkloadRepository) CountByConfigID(configID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.K8sWorkload{}).Where("config_id = ?", configID).Count(&count).Error
	return count, err
}

// BatchUpdate 批量更新工作负载
func (r *k8sWorkloadRepository) BatchUpdate(workloads []model.K8sWorkload) error {
	if len(workloads) == 0 {
		return nil
	}
	
	// 使用事务批量更新
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, workload := range workloads {
			if err := tx.Save(&workload).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
