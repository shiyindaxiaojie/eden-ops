package repository

import (
	"eden-ops/internal/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// K8sWorkloadRepository Kubernetes工作负载仓库接口
type K8sWorkloadRepository interface {
	Create(workload *model.K8sWorkload) error
	Update(workload *model.K8sWorkload) error
	Delete(id int64) error
	Get(id int64) (*model.K8sWorkload, error)
	List(configID int64, page, pageSize int) (int64, []model.K8sWorkload, error)
	ListWithFilter(page, pageSize int, name, namespace, workloadType, status, replicas, sortBy, sortOrder string, startTime, endTime *string, configId *int64) (int64, []model.K8sWorkload, error)
	ListByConfigID(configID int64) ([]model.K8sWorkload, error)
	DeleteByConfigID(configID int64) error
	DeleteNotInList(configID int64, currentWorkloads []model.K8sWorkload) error
	BatchCreate(workloads []model.K8sWorkload) error
	BatchUpdate(workloads []model.K8sWorkload) error
	BatchCreateOrUpdate(workloads []model.K8sWorkload) error
	CountByConfigID(configID int64) (int64, error)
	// 事务支持
	WithTx(tx *gorm.DB) K8sWorkloadRepository
	Transaction(fn func(K8sWorkloadRepository) error) error
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
func (r *k8sWorkloadRepository) ListWithFilter(page, pageSize int, name, namespace, workloadType, status, replicas, sortBy, sortOrder string, startTime, endTime *string, configId *int64) (int64, []model.K8sWorkload, error) {
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
	if status != "" {
		if status == "Running" {
			// 运行中：replicas 大于 0
			query = query.Where("replicas > 0")
		} else if status == "Other" {
			// 闲置：replicas 等于 0
			query = query.Where("replicas = 0")
		} else {
			// 其他状态按原来的逻辑
			query = query.Where("status = ?", status)
		}
	}

	// replicas过滤
	if replicas != "" {
		if replicas == "gt0" {
			// 大于0
			query = query.Where("replicas > 0")
		} else if replicas == "eq0" {
			// 等于0
			query = query.Where("replicas = 0")
		}
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
	if err := query.Offset(offset).Limit(pageSize).Order(orderClause).Find(&workloads).Error; err != nil {
		return 0, nil, err
	}

	return total, workloads, nil
}

// buildOrderClause 构建排序条件
func (r *k8sWorkloadRepository) buildOrderClause(sortBy, sortOrder string) string {
	// 默认排序：状态优先级 + 创建时间倒序
	defaultOrder := `
		CASE status
			WHEN 'Pending' THEN 1
			WHEN 'Failed' THEN 2
			WHEN 'Error' THEN 3
			WHEN 'Progressing' THEN 4
			WHEN 'Running' THEN 5
			WHEN 'Available' THEN 6
			WHEN 'Complete' THEN 7
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
	case "kind":
		return "kind " + sortOrder + ", " + defaultOrder
	case "status":
		if sortOrder == "asc" {
			return defaultOrder
		} else {
			return `
				CASE status
					WHEN 'Complete' THEN 1
					WHEN 'Available' THEN 2
					WHEN 'Running' THEN 3
					WHEN 'Progressing' THEN 4
					WHEN 'Error' THEN 5
					WHEN 'Failed' THEN 6
					WHEN 'Pending' THEN 7
					ELSE 8
				END ASC, created_at DESC`
		}
	case "replicas":
		return "replicas " + sortOrder + ", " + defaultOrder
	case "ready_replicas":
		return "ready_replicas " + sortOrder + ", " + defaultOrder
	case "created_at":
		return "created_at " + sortOrder + ", name ASC"
	case "updated_at":
		return "updated_at " + sortOrder + ", name ASC"
	default:
		return defaultOrder
	}
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

// DeleteNotInList 删除不在当前列表中的工作负载
func (r *k8sWorkloadRepository) DeleteNotInList(configID int64, currentWorkloads []model.K8sWorkload) error {
	if len(currentWorkloads) == 0 {
		// 如果没有当前工作负载，删除所有工作负载
		return r.db.Where("config_id = ? AND deleted_at IS NULL", configID).Delete(&model.K8sWorkload{}).Error
	}

	// 构建当前工作负载的唯一标识列表 (name-namespace-kind)
	var currentKeys []string
	for _, workload := range currentWorkloads {
		key := fmt.Sprintf("'%s-%s-%s'",
			strings.ReplaceAll(workload.Name, "'", "''"),
			strings.ReplaceAll(workload.Namespace, "'", "''"),
			strings.ReplaceAll(workload.Kind, "'", "''"))
		currentKeys = append(currentKeys, key)
	}

	// 删除不在当前列表中的工作负载
	sql := `DELETE FROM infra_k8s_workload
			WHERE config_id = ? AND deleted_at IS NULL
			AND CONCAT(name, '-', namespace, '-', kind) NOT IN (` + strings.Join(currentKeys, ",") + `)`

	return r.db.Exec(sql, configID).Error
}

// BatchCreateOrUpdate 批量创建或更新工作负载
func (r *k8sWorkloadRepository) BatchCreateOrUpdate(workloads []model.K8sWorkload) error {
	if len(workloads) == 0 {
		return nil
	}

	// 分批处理，减少锁持有时间
	batchSize := 50 // 每批处理50个
	for i := 0; i < len(workloads); i += batchSize {
		end := i + batchSize
		if end > len(workloads) {
			end = len(workloads)
		}

		batch := workloads[i:end]

		// 使用ON DUPLICATE KEY UPDATE避免竞态条件
		err := r.db.Transaction(func(tx *gorm.DB) error {
			// 使用原生SQL进行批量插入或更新，避免唯一键冲突
			return r.batchUpsertWorkloads(tx, batch)
		})

		if err != nil {
			return fmt.Errorf("batch %d-%d failed: %v", i, end-1, err)
		}
	}

	return nil
}

// batchUpsertWorkloads 批量插入或更新工作负载（使用原生SQL）
func (r *k8sWorkloadRepository) batchUpsertWorkloads(tx *gorm.DB, workloads []model.K8sWorkload) error {
	if len(workloads) == 0 {
		return nil
	}

	// 构建批量插入SQL
	sql := `INSERT INTO infra_k8s_workload
		(config_id, name, namespace, kind, replicas, ready_replicas, status, labels, selector, images,
		 cpu_request, cpu_limit, memory_request, memory_limit, created_at, updated_at)
		VALUES `

	var values []string
	var args []interface{}

	for _, workload := range workloads {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			workload.ConfigID, workload.Name, workload.Namespace, workload.Kind,
			workload.Replicas, workload.ReadyReplicas, workload.Status,
			workload.Labels, workload.Selector, workload.Images,
			workload.CPURequest, workload.CPULimit, workload.MemoryRequest, workload.MemoryLimit,
			workload.CreatedAt, workload.UpdatedAt)
	}

	sql += strings.Join(values, ", ")
	sql += ` ON DUPLICATE KEY UPDATE
		replicas = VALUES(replicas),
		ready_replicas = VALUES(ready_replicas),
		status = VALUES(status),
		labels = VALUES(labels),
		selector = VALUES(selector),
		images = VALUES(images),
		cpu_request = VALUES(cpu_request),
		cpu_limit = VALUES(cpu_limit),
		memory_request = VALUES(memory_request),
		memory_limit = VALUES(memory_limit),
		updated_at = VALUES(updated_at)`

	return tx.Exec(sql, args...).Error
}

// WithTx 使用事务
func (r *k8sWorkloadRepository) WithTx(tx *gorm.DB) K8sWorkloadRepository {
	return &k8sWorkloadRepository{db: tx}
}

// Transaction 执行事务
func (r *k8sWorkloadRepository) Transaction(fn func(K8sWorkloadRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTx(tx))
	})
}
