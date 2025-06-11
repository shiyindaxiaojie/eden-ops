package service

import (
	"context"
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sConfigService Kubernetes配置服务接口
type K8sConfigService interface {
	Create(config *model.K8sConfig) error
	Update(config *model.K8sConfig) error
	Delete(id uint) error
	Get(id uint) (*model.K8sConfig, error)
	List(page, pageSize int) ([]*model.K8sConfig, int64, error)
	TestConnection(config *model.K8sConfig) error
	SyncCluster(id int64) error
	GetNamespaces(id int64) ([]string, error)
}

// k8sConfigService Kubernetes配置服务实现
type k8sConfigService struct {
	repo            repository.K8sConfigRepository
	workloadService K8sWorkloadService
}

// NewK8sConfigService 创建Kubernetes配置服务
func NewK8sConfigService(repo repository.K8sConfigRepository, workloadService K8sWorkloadService) K8sConfigService {
	return &k8sConfigService{
		repo:            repo,
		workloadService: workloadService,
	}
}

// Create 创建Kubernetes配置
func (s *k8sConfigService) Create(config *model.K8sConfig) error {
	return s.repo.Create(config)
}

// Update 更新Kubernetes配置
func (s *k8sConfigService) Update(config *model.K8sConfig) error {
	return s.repo.Update(config)
}

// Delete 删除Kubernetes配置
func (s *k8sConfigService) Delete(id uint) error {
	return s.repo.Delete(int64(id))
}

// Get 获取Kubernetes配置
func (s *k8sConfigService) Get(id uint) (*model.K8sConfig, error) {
	return s.repo.Get(int64(id))
}

// List 获取Kubernetes配置列表
func (s *k8sConfigService) List(page, pageSize int) ([]*model.K8sConfig, int64, error) {
	total, configs, err := s.repo.List(page, pageSize, "")
	if err != nil {
		return nil, 0, err
	}

	var result []*model.K8sConfig
	for i := range configs {
		result = append(result, &configs[i])
	}

	return result, total, nil
}

// TestConnection 测试Kubernetes连接
func (s *k8sConfigService) TestConnection(config *model.K8sConfig) error {
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
		return fmt.Errorf("failed to get cluster version: %v", err)
	}

	return nil
}

// SyncCluster 同步集群信息
func (s *k8sConfigService) SyncCluster(id int64) error {
	// 获取集群配置
	config, err := s.repo.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get kubernetes config: %v", err)
	}

	// 检查集群是否启用
	if config.Status != 1 {
		return fmt.Errorf("cluster %s is disabled", config.Name)
	}

	// 解析kubeconfig
	k8sConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(config.Kubeconfig))
	if err != nil {
		return fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	// 获取集群版本
	version, err := clientset.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get cluster version: %v", err)
	}

	// 更新集群版本和最后同步时间
	now := time.Now()
	config.Version = version.String()
	config.LastSyncTime = &now
	if err := s.repo.Update(config); err != nil {
		return fmt.Errorf("failed to update cluster info: %v", err)
	}

	// 获取工作负载
	workloads, err := s.getWorkloadsFromCluster(clientset, id)
	if err != nil {
		return fmt.Errorf("failed to get workloads: %v", err)
	}

	// 同步工作负载到数据库
	if err := s.workloadService.SyncWorkloads(id, workloads); err != nil {
		return fmt.Errorf("failed to sync workloads: %v", err)
	}

	return nil
}

// GetNamespaces 获取命名空间列表
func (s *k8sConfigService) GetNamespaces(id int64) ([]string, error) {
	// 简化实现，返回默认命名空间
	return []string{"default", "kube-system", "kube-public"}, nil
}

// getWorkloadsFromCluster 从K8s集群获取工作负载
func (s *k8sConfigService) getWorkloadsFromCluster(clientset *kubernetes.Clientset, configID int64) ([]model.K8sWorkload, error) {
	var workloads []model.K8sWorkload
	ctx := context.Background()

	// 获取Deployments
	deployments, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}
	for _, d := range deployments.Items {
		replicas := int32(0)
		if d.Spec.Replicas != nil {
			replicas = *d.Spec.Replicas
		}
		workload := model.K8sWorkload{
			ConfigID:      configID,
			Name:          d.Name,
			Namespace:     d.Namespace,
			Kind:          "Deployment",
			Replicas:      int(replicas),
			ReadyReplicas: int(d.Status.ReadyReplicas),
			Status:        getDeploymentStatus(d.Status.Conditions),
		}
		workloads = append(workloads, workload)
	}

	// 获取StatefulSets
	statefulsets, err := clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list statefulsets: %v", err)
	}
	for _, s := range statefulsets.Items {
		replicas := int32(0)
		if s.Spec.Replicas != nil {
			replicas = *s.Spec.Replicas
		}
		workload := model.K8sWorkload{
			ConfigID:      configID,
			Name:          s.Name,
			Namespace:     s.Namespace,
			Kind:          "StatefulSet",
			Replicas:      int(replicas),
			ReadyReplicas: int(s.Status.ReadyReplicas),
			Status:        getStatefulSetStatus(s.Status.Conditions),
		}
		workloads = append(workloads, workload)
	}

	// 获取DaemonSets
	daemonsets, err := clientset.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list daemonsets: %v", err)
	}
	for _, d := range daemonsets.Items {
		workload := model.K8sWorkload{
			ConfigID:      configID,
			Name:          d.Name,
			Namespace:     d.Namespace,
			Kind:          "DaemonSet",
			Replicas:      int(d.Status.DesiredNumberScheduled),
			ReadyReplicas: int(d.Status.NumberReady),
			Status:        getDaemonSetStatus(d.Status.Conditions),
		}
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// getDeploymentStatus 获取Deployment状态
func getDeploymentStatus(conditions []metav1.Condition) string {
	for _, condition := range conditions {
		if condition.Type == "Available" && condition.Status == "True" {
			return "Running"
		}
	}
	return "Pending"
}

// getStatefulSetStatus 获取StatefulSet状态
func getStatefulSetStatus(conditions []metav1.Condition) string {
	for _, condition := range conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return "Running"
		}
	}
	return "Pending"
}

// getDaemonSetStatus 获取DaemonSet状态
func getDaemonSetStatus(conditions []metav1.Condition) string {
	for _, condition := range conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return "Running"
		}
	}
	return "Pending"
}
