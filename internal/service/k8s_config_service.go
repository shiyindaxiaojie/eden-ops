package service

import (
	"context"
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sConfigService Kubernetes配置服务接口
type K8sConfigService interface {
	Create(config *model.K8sConfig) error
	CreateWithClusterInfo(config *model.K8sConfig) error
	Update(config *model.K8sConfig) error
	UpdateWithClusterInfo(config *model.K8sConfig) error
	Delete(id uint) error
	Get(id uint) (*model.K8sConfig, error)
	List(page, pageSize int, name string, status *int, providerId *int64) ([]*model.K8sConfig, int64, error)
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

// CreateWithClusterInfo 创建Kubernetes配置并获取集群信息
func (s *k8sConfigService) CreateWithClusterInfo(config *model.K8sConfig) error {
	// 先创建配置
	if err := s.repo.Create(config); err != nil {
		return err
	}

	// 如果启用状态，则获取集群信息
	if config.Status == 1 {
		if err := s.updateClusterInfo(config); err != nil {
			// 记录错误但不影响创建
			fmt.Printf("Warning: failed to get cluster info: %v\n", err)
		}
	}

	return nil
}

// Update 更新Kubernetes配置
func (s *k8sConfigService) Update(config *model.K8sConfig) error {
	return s.repo.Update(config)
}

// UpdateWithClusterInfo 更新Kubernetes配置并获取集群信息
func (s *k8sConfigService) UpdateWithClusterInfo(config *model.K8sConfig) error {
	// 先更新配置
	if err := s.repo.Update(config); err != nil {
		return err
	}

	// 如果启用状态，则获取集群信息
	if config.Status == 1 {
		if err := s.updateClusterInfo(config); err != nil {
			// 记录错误但不影响更新
			fmt.Printf("Warning: failed to get cluster info: %v\n", err)
		}
	}

	return nil
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
func (s *k8sConfigService) List(page, pageSize int, name string, status *int, providerId *int64) ([]*model.K8sConfig, int64, error) {
	total, configs, err := s.repo.List(page, pageSize, name, status, providerId)
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

// updateClusterInfo 更新集群信息（不同步工作负载）
func (s *k8sConfigService) updateClusterInfo(config *model.K8sConfig) error {
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

	// 获取集群资源信息
	clusterInfo, err := s.getClusterInfo(clientset, k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to get cluster info: %v", err)
	}

	// 更新集群信息
	now := time.Now()
	config.Version = version.String()
	config.Context = clusterInfo.Context
	config.NodeCount = clusterInfo.NodeCount
	config.PodCount = clusterInfo.PodCount
	config.CPUTotal = clusterInfo.CPUTotal
	config.CPUUsed = clusterInfo.CPUUsed
	config.MemoryTotal = clusterInfo.MemoryTotal
	config.MemoryUsed = clusterInfo.MemoryUsed
	config.LastSyncTime = &now

	// 更新到数据库
	return s.repo.Update(config)
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

	// 获取集群资源信息
	clusterInfo, err := s.getClusterInfo(clientset, k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to get cluster info: %v", err)
	}

	// 更新集群信息
	now := time.Now()
	config.Version = version.String()
	config.Context = clusterInfo.Context
	config.NodeCount = clusterInfo.NodeCount
	config.PodCount = clusterInfo.PodCount
	config.CPUTotal = clusterInfo.CPUTotal
	config.CPUUsed = clusterInfo.CPUUsed
	config.MemoryTotal = clusterInfo.MemoryTotal
	config.MemoryUsed = clusterInfo.MemoryUsed
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
			Status:        getDeploymentStatus(d.Status),
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
			Status:        getStatefulSetStatus(s.Status),
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
			Status:        getDaemonSetStatus(d.Status),
		}
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// ClusterInfo 集群信息结构
type ClusterInfo struct {
	Context     string
	NodeCount   int
	PodCount    int
	CPUTotal    string
	CPUUsed     string
	MemoryTotal string
	MemoryUsed  string
}

// getClusterInfo 获取集群资源信息
func (s *k8sConfigService) getClusterInfo(clientset *kubernetes.Clientset, k8sConfig *rest.Config) (*ClusterInfo, error) {
	ctx := context.Background()
	info := &ClusterInfo{}

	// 获取节点信息
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}
	info.NodeCount = len(nodes.Items)

	// 计算总CPU和内存
	var totalCPU, totalMemory resource.Quantity
	for _, node := range nodes.Items {
		if cpu, ok := node.Status.Capacity[corev1.ResourceCPU]; ok {
			totalCPU.Add(cpu)
		}
		if memory, ok := node.Status.Capacity[corev1.ResourceMemory]; ok {
			totalMemory.Add(memory)
		}
	}
	info.CPUTotal = totalCPU.String()
	info.MemoryTotal = formatMemory(totalMemory.Value())

	// 获取Pod信息
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}
	info.PodCount = len(pods.Items)

	// 计算已使用的CPU和内存（简化计算，基于Pod请求）
	var usedCPU, usedMemory resource.Quantity
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if cpu, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				usedCPU.Add(cpu)
			}
			if memory, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				usedMemory.Add(memory)
			}
		}
	}
	info.CPUUsed = usedCPU.String()
	info.MemoryUsed = formatMemory(usedMemory.Value())

	// 获取Context信息（从kubeconfig中解析）
	info.Context = "default" // 简化实现

	return info, nil
}

// formatMemory 格式化内存显示
func formatMemory(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// getDeploymentStatus 获取Deployment状态
func getDeploymentStatus(status appsv1.DeploymentStatus) string {
	for _, condition := range status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == "True" {
			return "Running"
		}
	}
	if status.ReadyReplicas > 0 {
		return "Running"
	}
	return "Pending"
}

// getStatefulSetStatus 获取StatefulSet状态
func getStatefulSetStatus(status appsv1.StatefulSetStatus) string {
	if status.ReadyReplicas > 0 {
		return "Running"
	}
	return "Pending"
}

// getDaemonSetStatus 获取DaemonSet状态
func getDaemonSetStatus(status appsv1.DaemonSetStatus) string {
	if status.NumberReady > 0 {
		return "Running"
	}
	return "Pending"
}
