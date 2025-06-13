package repository

import (
	"context"
	"eden-ops/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sConfigRepository Kubernetes配置仓库接口
type K8sConfigRepository interface {
	Create(config *model.K8sConfig) error
	Update(config *model.K8sConfig) error
	Delete(id int64) error
	Get(id int64) (*model.K8sConfig, error)
	List(page, pageSize int, name string, status *int, providerId *int64, clusterID string) (int64, []model.K8sConfig, error)
	GetDB() *gorm.DB
}

// k8sConfigRepository Kubernetes配置仓库实现
type k8sConfigRepository struct {
	db *gorm.DB
}

// NewK8sConfigRepository 创建Kubernetes配置仓库实例
func NewK8sConfigRepository(db *gorm.DB) K8sConfigRepository {
	return &k8sConfigRepository{db: db}
}

// Create 创建Kubernetes配置
func (r *k8sConfigRepository) Create(config *model.K8sConfig) error {
	return r.db.Create(config).Error
}

// Update 更新Kubernetes配置
func (r *k8sConfigRepository) Update(config *model.K8sConfig) error {
	return r.db.Save(config).Error
}

// Delete 删除Kubernetes配置
func (r *k8sConfigRepository) Delete(id int64) error {
	result := r.db.Delete(&model.K8sConfig{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除kubernetes配置失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("kubernetes配置不存在: %d", id)
	}
	return nil
}

// Get 获取Kubernetes配置
func (r *k8sConfigRepository) Get(id int64) (*model.K8sConfig, error) {
	var config model.K8sConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("kubernetes配置不存在: %d", id)
		}
		return nil, fmt.Errorf("获取kubernetes配置失败: %v", err)
	}
	return &config, nil
}

// List 获取Kubernetes配置列表
func (r *k8sConfigRepository) List(page, pageSize int, name string, status *int, providerId *int64, clusterID string) (int64, []model.K8sConfig, error) {
	var configs []model.K8sConfig
	var total int64

	query := r.db.Model(&model.K8sConfig{})
	if name != "" {
		query = query.Where("infra_k8s_config.name LIKE ?", "%"+name+"%")
	}
	if status != nil {
		query = query.Where("infra_k8s_config.status = ?", *status)
	}
	if providerId != nil {
		query = query.Where("infra_k8s_config.provider_id = ?", *providerId)
	}
	if clusterID != "" {
		query = query.Where("infra_k8s_config.cluster_id = ?", clusterID)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Select("infra_k8s_config.*, infra_cloud_provider.name as provider_name").
		Joins("LEFT JOIN infra_cloud_provider ON infra_k8s_config.provider_id = infra_cloud_provider.id").
		Scan(&configs).Error; err != nil {
		return 0, nil, err
	}

	// 手动设置ProviderName字段
	for i := range configs {
		if configs[i].ProviderId != nil {
			var providerName string
			if err := r.db.Model(&model.CloudProvider{}).
				Where("id = ?", *configs[i].ProviderId).
				Select("name").
				Scan(&providerName).Error; err == nil {
				configs[i].ProviderName = providerName
			}
		}
	}

	return total, configs, nil
}

// GetDB 获取数据库连接
func (r *k8sConfigRepository) GetDB() *gorm.DB {
	return r.db
}

// TestConnection 测试Kubernetes集群连接
func (r *k8sConfigRepository) TestConnection(config *model.K8sConfig) error {
	k8sConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(config.Kubeconfig))
	if err != nil {
		return fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to connect to kubernetes cluster: %v", err)
	}

	return nil
}

// GetWorkloads 获取Kubernetes集群工作负载
func (r *k8sConfigRepository) GetWorkloads(id int64) (interface{}, error) {
	config, err := r.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %v", err)
	}

	k8sConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(config.Kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	// Get deployments from all namespaces
	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}

	var workloads []model.K8sWorkloadInfo
	for _, deployment := range deployments.Items {
		workload := model.K8sWorkloadInfo{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Kind:      "Deployment",
			Replicas:  int(*deployment.Spec.Replicas),
			Status:    string(deployment.Status.Conditions[0].Type),
		}
		workloads = append(workloads, workload)
	}

	return workloads, nil
}
