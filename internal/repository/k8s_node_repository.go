package repository

import (
	"eden-ops/internal/model"
	"gorm.io/gorm"
)

// K8sNodeRepository 节点仓储接口
type K8sNodeRepository interface {
	Create(node *model.K8sNode) error
	Update(node *model.K8sNode) error
	Delete(id int64) error
	GetByID(id int64) (*model.K8sNode, error)
	GetByConfigAndName(configID int64, name string) (*model.K8sNode, error)
	List(page, pageSize int, configID int64, name, internalIP, status string, ready *bool) (int64, []model.K8sNode, error)
	DeleteByConfigID(configID int64) error
	BatchCreateOrUpdate(nodes []model.K8sNode) error
}

// k8sNodeRepository 节点仓储实现
type k8sNodeRepository struct {
	db *gorm.DB
}

// NewK8sNodeRepository 创建节点仓储
func NewK8sNodeRepository(db *gorm.DB) K8sNodeRepository {
	return &k8sNodeRepository{db: db}
}

// Create 创建节点
func (r *k8sNodeRepository) Create(node *model.K8sNode) error {
	return r.db.Create(node).Error
}

// Update 更新节点
func (r *k8sNodeRepository) Update(node *model.K8sNode) error {
	return r.db.Save(node).Error
}

// Delete 删除节点
func (r *k8sNodeRepository) Delete(id int64) error {
	return r.db.Delete(&model.K8sNode{}, id).Error
}

// GetByID 根据ID获取节点
func (r *k8sNodeRepository) GetByID(id int64) (*model.K8sNode, error) {
	var node model.K8sNode
	err := r.db.First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetByConfigAndName 根据配置ID和名称获取节点
func (r *k8sNodeRepository) GetByConfigAndName(configID int64, name string) (*model.K8sNode, error) {
	var node model.K8sNode
	err := r.db.Where("config_id = ? AND name = ?", configID, name).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// List 获取节点列表
func (r *k8sNodeRepository) List(page, pageSize int, configID int64, name, internalIP, status string, ready *bool) (int64, []model.K8sNode, error) {
	var nodes []model.K8sNode
	var total int64

	query := r.db.Model(&model.K8sNode{})

	if configID > 0 {
		query = query.Where("config_id = ?", configID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if internalIP != "" {
		query = query.Where("internal_ip LIKE ?", "%"+internalIP+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if ready != nil {
		query = query.Where("ready = ?", *ready)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&nodes).Error; err != nil {
		return 0, nil, err
	}

	return total, nodes, nil
}

// DeleteByConfigID 根据配置ID删除所有节点
func (r *k8sNodeRepository) DeleteByConfigID(configID int64) error {
	return r.db.Where("config_id = ?", configID).Delete(&model.K8sNode{}).Error
}

// BatchCreateOrUpdate 批量创建或更新节点
func (r *k8sNodeRepository) BatchCreateOrUpdate(nodes []model.K8sNode) error {
	if len(nodes) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, node := range nodes {
			var existingNode model.K8sNode
			err := tx.Where("config_id = ? AND name = ?", node.ConfigID, node.Name).First(&existingNode).Error
			
			if err == gorm.ErrRecordNotFound {
				// 创建新节点
				if err := tx.Create(&node).Error; err != nil {
					return err
				}
			} else if err == nil {
				// 更新现有节点
				node.ID = existingNode.ID
				node.CreatedAt = existingNode.CreatedAt
				if err := tx.Save(&node).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})
}
