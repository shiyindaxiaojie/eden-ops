package repository

import (
	"eden-ops/internal/model"
	"fmt"
	"strings"

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
	DeleteNotInList(configID int64, currentPods []model.K8sPod) error
	BatchCreate(pods []model.K8sPod) error
	BatchCreateOrUpdate(pods []model.K8sPod) error
	// 事务支持
	WithTx(tx *gorm.DB) K8sPodRepository
	Transaction(fn func(K8sPodRepository) error) error
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
		if status == "Error" {
			// 异常状态：非Running状态
			query = query.Where("status != ?", "Running")
		} else {
			// 其他状态按原来的逻辑
			query = query.Where("status = ?", status)
		}
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

// DeleteNotInList 删除不在当前列表中的Pod
func (r *k8sPodRepository) DeleteNotInList(configID int64, currentPods []model.K8sPod) error {
	if len(currentPods) == 0 {
		// 如果没有当前Pod，删除所有Pod
		return r.db.Where("config_id = ? AND deleted_at IS NULL", configID).Delete(&model.K8sPod{}).Error
	}

	// 构建当前Pod的唯一标识列表 (name-namespace)
	var currentKeys []string
	for _, pod := range currentPods {
		key := fmt.Sprintf("'%s-%s'",
			strings.ReplaceAll(pod.Name, "'", "''"),
			strings.ReplaceAll(pod.Namespace, "'", "''"))
		currentKeys = append(currentKeys, key)
	}

	// 删除不在当前列表中的Pod
	sql := `DELETE FROM infra_k8s_pod
			WHERE config_id = ? AND deleted_at IS NULL
			AND CONCAT(name, '-', namespace) NOT IN (` + strings.Join(currentKeys, ",") + `)`

	return r.db.Exec(sql, configID).Error
}

// BatchCreateOrUpdate 批量创建或更新Pod
func (r *k8sPodRepository) BatchCreateOrUpdate(pods []model.K8sPod) error {
	if len(pods) == 0 {
		return nil
	}

	// 分批处理，减少锁持有时间
	batchSize := 50 // 每批处理50个
	for i := 0; i < len(pods); i += batchSize {
		end := i + batchSize
		if end > len(pods) {
			end = len(pods)
		}

		batch := pods[i:end]

		// 使用较短的事务处理每批数据
		err := r.db.Transaction(func(tx *gorm.DB) error {
			for _, pod := range batch {
				// 先查找是否存在相同的Pod
				var existingPod model.K8sPod
				err := tx.Where("config_id = ? AND name = ? AND namespace = ?",
					pod.ConfigID, pod.Name, pod.Namespace).First(&existingPod).Error

				if err == gorm.ErrRecordNotFound {
					// 不存在，创建新记录
					if err := tx.Create(&pod).Error; err != nil {
						return err
					}
				} else if err == nil {
					// 存在，更新记录
					pod.ID = existingPod.ID
					pod.CreatedAt = existingPod.CreatedAt
					if err := tx.Save(&pod).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("batch %d-%d failed: %v", i, end-1, err)
		}
	}

	return nil
}

// WithTx 使用事务
func (r *k8sPodRepository) WithTx(tx *gorm.DB) K8sPodRepository {
	return &k8sPodRepository{db: tx}
}

// Transaction 执行事务
func (r *k8sPodRepository) Transaction(fn func(K8sPodRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTx(tx))
	})
}
