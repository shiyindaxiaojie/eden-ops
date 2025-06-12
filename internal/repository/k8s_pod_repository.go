package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// K8sPodRepository K8s Pod仓库接口
type K8sPodRepository interface {
	Create(pod *model.K8sPod) error
	Update(pod *model.K8sPod) error
	Delete(id int64) error
	Get(id int64) (*model.K8sPod, error)
	List(configID int64, page, pageSize int) (int64, []model.K8sPod, error)
	ListWithFilter(page, pageSize int, name, namespace, workloadName, status, instanceIP, sortBy, sortOrder string, startTime, endTime *string, configId *int64) (int64, []model.K8sPod, error)
	ListByConfigID(configID int64) ([]model.K8sPod, error)
	DeleteByConfigID(configID int64) error
	BatchCreate(pods []model.K8sPod) error
}

// k8sPodRepository K8s Pod仓库实现
type k8sPodRepository struct {
	db *gorm.DB
}

// NewK8sPodRepository 创建K8s Pod仓库
func NewK8sPodRepository(db *gorm.DB) K8sPodRepository {
	return &k8sPodRepository{db: db}
}

// Create 创建Pod
func (r *k8sPodRepository) Create(pod *model.K8sPod) error {
	return r.db.Create(pod).Error
}

// Update 更新Pod
func (r *k8sPodRepository) Update(pod *model.K8sPod) error {
	return r.db.Save(pod).Error
}

// Delete 删除Pod
func (r *k8sPodRepository) Delete(id int64) error {
	return r.db.Delete(&model.K8sPod{}, id).Error
}

// Get 获取Pod
func (r *k8sPodRepository) Get(id int64) (*model.K8sPod, error) {
	var pod model.K8sPod
	err := r.db.First(&pod, id).Error
	if err != nil {
		return nil, err
	}
	return &pod, nil
}

// List 获取Pod列表
func (r *k8sPodRepository) List(configID int64, page, pageSize int) (int64, []model.K8sPod, error) {
	var pods []model.K8sPod
	var total int64

	query := r.db.Model(&model.K8sPod{}).Where("config_id = ?", configID)

	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&pods).Error; err != nil {
		return 0, nil, err
	}

	return total, pods, nil
}

// ListWithFilter 获取Pod列表（支持筛选）
func (r *k8sPodRepository) ListWithFilter(page, pageSize int, name, namespace, workloadName, status, instanceIP, sortBy, sortOrder string, startTime, endTime *string, configId *int64) (int64, []model.K8sPod, error) {
	var pods []model.K8sPod
	var total int64

	query := r.db.Model(&model.K8sPod{})

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
	if workloadName != "" {
		query = query.Where("workload_name LIKE ?", "%"+workloadName+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if instanceIP != "" {
		query = query.Where("instance_ip LIKE ?", "%"+instanceIP+"%")
	}

	// 时间范围筛选
	if startTime != nil && *startTime != "" {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil && *endTime != "" {
		query = query.Where("created_at <= ?", *endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 构建排序条件
	orderClause := r.buildOrderClause(sortBy, sortOrder)
	
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order(orderClause).Find(&pods).Error; err != nil {
		return 0, nil, err
	}

	return total, pods, nil
}

// buildOrderClause 构建排序条件
func (r *k8sPodRepository) buildOrderClause(sortBy, sortOrder string) string {
	// 默认排序：状态优先级 + 创建时间倒序
	defaultOrder := `
		CASE status 
			WHEN 'Pending' THEN 1
			WHEN 'Failed' THEN 2
			WHEN 'Error' THEN 3
			WHEN 'CrashLoopBackOff' THEN 4
			WHEN 'ImagePullBackOff' THEN 5
			WHEN 'Running' THEN 6
			WHEN 'Succeeded' THEN 7
			ELSE 8
		END ASC, created_at DESC`
	
	if sortBy == "" {
		return defaultOrder
	}
	
	// 验证排序方向
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	
	// 根据排序字段构建排序条件
	switch sortBy {
	case "name":
		return "name " + sortOrder + ", " + defaultOrder
	case "namespace":
		return "namespace " + sortOrder + ", " + defaultOrder
	case "workload_name":
		return "workload_name " + sortOrder + ", " + defaultOrder
	case "workload_kind":
		return "workload_kind " + sortOrder + ", " + defaultOrder
	case "status":
		if sortOrder == "asc" {
			return defaultOrder
		} else {
			return `
				CASE status 
					WHEN 'Succeeded' THEN 1
					WHEN 'Running' THEN 2
					WHEN 'ImagePullBackOff' THEN 3
					WHEN 'CrashLoopBackOff' THEN 4
					WHEN 'Error' THEN 5
					WHEN 'Failed' THEN 6
					WHEN 'Pending' THEN 7
					ELSE 8
				END ASC, created_at DESC`
		}
	case "restart_count":
		return "restart_count " + sortOrder + ", " + defaultOrder
	case "start_time":
		return "start_time " + sortOrder + ", name ASC"
	case "created_at":
		return "created_at " + sortOrder + ", name ASC"
	default:
		return defaultOrder
	}
}

// ListByConfigID 根据配置ID获取所有Pod
func (r *k8sPodRepository) ListByConfigID(configID int64) ([]model.K8sPod, error) {
	var pods []model.K8sPod
	err := r.db.Where("config_id = ?", configID).Find(&pods).Error
	return pods, err
}

// DeleteByConfigID 根据配置ID删除所有Pod
func (r *k8sPodRepository) DeleteByConfigID(configID int64) error {
	return r.db.Where("config_id = ?", configID).Delete(&model.K8sPod{}).Error
}

// BatchCreate 批量创建Pod
func (r *k8sPodRepository) BatchCreate(pods []model.K8sPod) error {
	if len(pods) == 0 {
		return nil
	}
	return r.db.CreateInBatches(pods, 100).Error
}
